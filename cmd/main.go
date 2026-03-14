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
	cmds.Register("login", commands.Login)
	cmds.Register("register", commands.Register)
	cmds.Register("reset", commands.Reset)
	cmds.Register("users", commands.Users)
	cmds.Register("agg", commands.Aggregate)
	cmds.Register("addfeed", commands.AddFeed)
	cmds.Register("feeds", commands.Feeds)
	cmds.Register("follow", commands.Follow)
	cmds.Register("following", commands.Following)

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("usage: gator <command> [args]")
		os.Exit(1)
	}

	if err := cmds.Run(st, commands.Command{Name: args[0], Args: args[1:]}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
