// Package main is the entry point for gator, a CLI RSS feed aggregator.
package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ibrais/gator/internal/commands"
	"github.com/ibrais/gator/internal/config"
	database "github.com/ibrais/gator/internal/database/generated"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error loading config:", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Println("error connecting to DB:", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	st := &commands.State{
		Config: cfg,
		DB:     dbQueries,
	}

	cmds := commands.New()
	// User management
	cmds.Register("login", commands.Login)
	cmds.Register("register", commands.Register)
	cmds.Register("users", commands.Users)

	// Feed management
	cmds.Register("addfeed", commands.MiddlewareLoggedIn(commands.AddFeed))
	cmds.Register("feeds", commands.Feeds)
	cmds.Register("follow", commands.MiddlewareLoggedIn(commands.Follow))
	cmds.Register("following", commands.MiddlewareLoggedIn(commands.Following))
	cmds.Register("unfollow", commands.MiddlewareLoggedIn(commands.Unfollow))

	// Aggregation
	cmds.Register("agg", commands.MiddlewareLoggedIn(commands.Aggregate))
	cmds.Register("browse", commands.MiddlewareLoggedIn(commands.Browse))

	// Delete all users (dev/debug use)
	cmds.Register("reset", commands.Reset)

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("usage: gator <command> [args]")
		os.Exit(1)
	}

	if err := cmds.Run(st, commands.Command{
		Name: args[0],
		Args: args[1:]},
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
