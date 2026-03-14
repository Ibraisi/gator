// Package commands defines the CLI command registry, shared state, and all command handlers for gator.
package commands

import (
	"errors"

	"github.com/ibrais/gator/internal/config"
	database "github.com/ibrais/gator/internal/database/generated"
)

type State struct {
	Config *config.Config
	DB     *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	handlers map[string]func(*State, Command) error
}

func New() *Commands {
	return &Commands{
		handlers: make(map[string]func(*State, Command) error),
	}
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.handlers[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.handlers[cmd.Name]
	if !ok {
		return errors.New("command not found in registry")
	}
	return handler(s, cmd)
}
