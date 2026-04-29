package commands

import (
	"context"
	"fmt"

	"github.com/ahmed-abdelhamid/gator/internal/cli"
)

// Reset wipes all users from the database. Destructive; takes no arguments.
func Reset(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("reset takes no arguments")
	}

	if err := s.DB.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("delete users: %w", err)
	}

	fmt.Println("all users deleted")
	return nil
}
