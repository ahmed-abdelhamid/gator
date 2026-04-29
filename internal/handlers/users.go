package handlers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ahmed-abdelhamid/gator/internal/cli"
	"github.com/ahmed-abdelhamid/gator/internal/database"
	"github.com/google/uuid"
)

func Login(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login requires a username argument")
	}

	user, err := s.DB.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		os.Exit(1)
	}

	if err := s.Cfg.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("user has been set to %s\n", user.Name)
	return nil
}

func Register(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("register requires a username argument")
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}
	user, err := s.DB.CreateUser(context.Background(), params)
	if err != nil {
		os.Exit(1)
	}

	if err := s.Cfg.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("new user created: %s\n", user.Name)
	return nil
}

func Reset(s *cli.State, cmd cli.Command) error {
	err := s.DB.DeleteUsers(context.Background())
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("all users deleted\n")
	return nil
}

func Users(s *cli.State, cmd cli.Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		os.Exit(1)
	}

	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if s.Cfg.CurrentUserName == user.Name {
			fmt.Printf(" (current)")
		}
		fmt.Printf("\n")
	}

	return nil
}
