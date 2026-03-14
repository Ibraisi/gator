package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	database "github.com/ibrais/gator/internal/database/generated"
	"github.com/ibrais/gator/internal/rss"
)

func Login(s *State, cmd Command) error {
	if err := checkArgs(cmd.Args, 1, "login <username>"); err != nil {
		return err
	}
	name := cmd.Args[0]
	user, err := s.DB.GetUserByName(context.Background(), name)
	if err != nil {
		return err
	}

	if err := s.Config.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("User %s set successfully\n", name)
	return nil
}

func Register(s *State, cmd Command) error {
	if err := checkArgs(cmd.Args, 1, "register <username>"); err != nil {
		return err
	}
	user := cmd.Args[0]
	_, err := s.DB.CreateUser(context.Background(), user)
	if err != nil {
		return err
	}

	if err := s.Config.SetUser(user); err != nil {
		return err
	}

	fmt.Printf("User %s set successfully\n", user)
	return nil
}

func Reset(s *State, cmd Command) error {
	err := s.DB.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("All users deleted")
	return nil
}

func Users(s *State, cmd Command) error {
	if err := checkArgs(cmd.Args, 0, "users"); err != nil {
		return err
	}
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, u := range users {
		if u.Name == s.Config.CurrentUser {
			fmt.Printf("%s (current)\n", u.Name)
			continue
		}
		fmt.Println(u.Name)
	}
	return nil
}

func Aggregate(s *State, cmd Command, user database.User) error {
	if err := checkArgs(cmd.Args, 1, "agg <time_between_reqs>"); err != nil {
		return err
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	ctx := context.Background()
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		feeds, err := s.DB.GetNextFeedToFetch(ctx, user.ID)
		if err != nil {
			return err
		}
		if err := scrapeFeeds(ctx, s, feeds); err != nil {
			return err
		}
	}
}

func AddFeed(s *State, cmd Command, user database.User) error {
	if err := checkArgs(cmd.Args, 2, "addfeed <name> <url>"); err != nil {
		return err
	}
	ctx := context.Background()
	feedURL, err := url.Parse(cmd.Args[1])
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	feed, err := s.DB.CreateFeed(ctx, database.CreateFeedParams{
		Name:   cmd.Args[0],
		Url:    feedURL.String(),
		UserID: user.ID,
	})
	if err != nil {
		return err
	}

	_, err = s.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func Feeds(s *State, cmd Command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, f := range feeds {
		fmt.Println(f.FeedName)
		fmt.Println(f.Url)
		fmt.Println(f.UserName)
	}
	return nil
}

func Follow(s *State, cmd Command, user database.User) error {
	if err := checkArgs(cmd.Args, 1, "follow <url>"); err != nil {
		return err
	}
	feedURL := cmd.Args[0]
	ctx := context.Background()
	feed, err := s.DB.GetFeedByURL(ctx, feedURL)
	if err != nil {
		return err
	}

	dbRes, err := s.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println(dbRes.FeedName)
	fmt.Println(dbRes.UserName)

	return nil
}

func Following(s *State, cmd Command, user database.User) error {
	ctx := context.Background()
	rows, err := s.DB.GetFollowedFeedsNames(ctx, user.ID)
	if err != nil {
		return err
	}

	for _, r := range rows {
		fmt.Println(r.Name)
	}

	return nil
}

func Unfollow(s *State, cmd Command, user database.User) error {
	if err := checkArgs(cmd.Args, 1, "unfollow <url>"); err != nil {
		return err
	}
	return s.DB.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    cmd.Args[0],
	})
}

func Browse(s *State, cmd Command, user database.User) error {
	// Default 2 posts
	limit := int32(2)
	if len(cmd.Args) == 1 {
		n, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("usage: browse [limit]")
		}
		limit = int32(n)
	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: browse [limit]")
	}

	posts, err := s.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return err
	}

	for _, p := range posts {
		fmt.Printf("--- %s ---\n", p.Title)
		fmt.Printf("URL: %s\n", p.Url)
		if p.PublishedAt.Valid {
			fmt.Printf("Published: %s\n", p.PublishedAt.Time.Format(time.RFC1123))
		}
		if p.Description.Valid && p.Description.String != "" {
			fmt.Printf("%s\n", p.Description.String)
		}
		fmt.Println()
	}
	return nil
}

func scrapeFeeds(ctx context.Context, s *State, feeds []database.Feed) error {
	for _, f := range feeds {
		res, err := rss.FetchFeed(ctx, f.Url)
		if err != nil {
			return err
		}

		if err := s.DB.MarkFeedFetched(ctx, f.ID); err != nil {
			return err
		}

		for _, item := range res.Channel.Item {
			err := s.DB.CreatePost(ctx, database.CreatePostParams{
				Title:       item.Title,
				Url:         item.Link,
				Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
				PublishedAt: parsePubDate(item.PubDate),
				FeedID:      f.ID,
			})
			if err != nil {
				log.Printf("error saving post %q: %v", item.Link, err)
			}
		}
	}

	return nil
}

func parsePubDate(s string) sql.NullTime {
	if s == "" {
		return sql.NullTime{}
	}
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
	}
	for _, f := range formats {
		t, err := time.Parse(f, s)
		if err == nil {
			return sql.NullTime{Time: t, Valid: true}
		}
	}
	log.Printf("could not parse pub date: %q", s)
	return sql.NullTime{}
}
