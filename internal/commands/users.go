// Package commands implements the gator CLI command handlers.
package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ahmed-abdelhamid/gator/internal/cli"
	"github.com/ahmed-abdelhamid/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Login switches the active user in the config to an existing DB user.
func Login(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login requires a username argument")
	}

	user, err := s.DB.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user %q does not exist", cmd.Args[0])
		}
		return fmt.Errorf("get user %q: %w", cmd.Args[0], err)
	}

	if err := s.Cfg.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("user has been set to %s\n", user.Name)
	return nil
}

// Register creates a new user and sets them as the active user.
func Register(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("register requires a username argument")
	}

	now := time.Now().UTC()
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.Args[0],
	}
	user, err := s.DB.CreateUser(context.Background(), params)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgUniqueViolation {
			return fmt.Errorf("user %q already exists", cmd.Args[0])
		}
		return fmt.Errorf("create user %q: %w", cmd.Args[0], err)
	}

	if err := s.Cfg.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("new user created: %s\n", user.Name)
	return nil
}

// Users prints every registered user, marking the active one.
func Users(s *cli.State, cmd cli.Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("get users: %w", err)
	}

	for _, user := range users {
		marker := ""
		if s.Cfg.CurrentUserName == user.Name {
			marker = " (current)"
		}
		fmt.Printf("* %s%s\n", user.Name, marker)
	}

	return nil
}
