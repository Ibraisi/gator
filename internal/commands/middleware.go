package commands

import (
	"context"
	"fmt"

	database "github.com/ibrais/gator/internal/database/generated"
)

func MiddlewareLoggedIn(handler func(*State, Command, database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.DB.GetUserByName(context.Background(), s.Config.CurrentUser)
		if err != nil {
			return fmt.Errorf("not logged in: %w", err)
		}
		return handler(s, cmd, user)
	}
}

func checkArgs(args []string, n int, usage string) error {
	if len(args) != n {
		return fmt.Errorf("usage: %s", usage)
	}
	return nil
}
