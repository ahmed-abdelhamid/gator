// Package cli wires the gator command dispatcher: a registry of named
// commands plus the per-invocation State they receive.
package cli

import (
	"database/sql"
	"fmt"

	"github.com/ahmed-abdelhamid/gator/internal/config"
	"github.com/ahmed-abdelhamid/gator/internal/database"
)

// State is the shared context handed to every command handler.
// Conn is the raw *sql.DB; reach for it only when you need a transaction.
// DB is the sqlc-generated Queries for normal use.
type State struct {
	DB   *database.Queries
	Conn *sql.DB
	Cfg  *config.Config
}

// Command is a parsed CLI invocation: a name plus its positional arguments.
type Command struct {
	Name string
	Args []string
}

// HandlerFunc is the signature every command handler implements.
type HandlerFunc func(*State, Command) error

// Commands is the dispatcher: a name -> handler registry.
type Commands struct {
	handlers map[string]HandlerFunc
}

// NewCommands returns an empty dispatcher.
func NewCommands() *Commands {
	return &Commands{handlers: map[string]HandlerFunc{}}
}

// Register binds a name to a handler. It panics on duplicate registration
// since that is always a programming error at startup.
func (c *Commands) Register(name string, f HandlerFunc) {
	if _, exists := c.handlers[name]; exists {
		panic(fmt.Sprintf("cli: command %q already registered", name))
	}
	c.handlers[name] = f
}

// Run dispatches cmd to its registered handler, returning an error for
// unknown commands or whatever the handler itself returns.
func (c *Commands) Run(s *State, cmd Command) error {
	h, ok := c.handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return h(s, cmd)
}
