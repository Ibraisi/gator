package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"github.com/google/uuid"
	database "github.com/ibrais/gator/internal/database/generated"
	"github.com/ibrais/gator/internal/rss"
)

func Login(s *State, cmd Command) error {
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
	user := cmd.Args[0]
	if err := s.Config.SetUser(user); err != nil {
		return err
	}
	createdUser, err := s.DB.CreateUser(context.Background(), user)
	if err != nil {
		return err
	}

	fmt.Printf("User %s set successfully\n", user)
	log.Println(createdUser)
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
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, u := range users {
		if u.Name == s.Config.CurrentUser {
			fmt.Printf("%s (current)", u.Name)
			continue
		}
		fmt.Println(u.Name)
	}
	return nil
}

func Aggregate(s *State, cmd Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	fmt.Println(feed)
	return err
}

func AddFeed(s *State, cmd Command) error {
	ctx := context.Background()
	feedURL, err := url.Parse(cmd.Args[1])
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	user, err := s.DB.GetUserByName(ctx, s.Config.CurrentUser)
	if err != nil {
		return err
	}
	feed, err := s.DB.CreateFeed(ctx, database.CreateFeedParams{
		Name:   sql.NullString{String: cmd.Args[0], Valid: true},
		Url:    sql.NullString{String: feedURL.String(), Valid: true},
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
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
		fmt.Println(f.Name)
		fmt.Println(f.Url)
		fmt.Println(f.UserName)
	}
	return nil
}
