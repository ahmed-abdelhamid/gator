package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/ahmed-abdelhamid/gator/internal/cli"
)

func Reset(s *cli.State, cmd cli.Command) error {
	err := s.DB.DeleteUsers(context.Background())
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("all users deleted\n")
	return nil
}
