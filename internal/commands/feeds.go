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

// AddFeed creates a feed owned by the currently logged-in user.
func AddFeed(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	if s.Cfg.CurrentUserName == "" {
		return fmt.Errorf("no user logged in; run `gator login` first")
	}

	user, err := s.DB.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("logged-in user %q does not exist", s.Cfg.CurrentUserName)
		}
		return fmt.Errorf("get user %q: %w", s.Cfg.CurrentUserName, err)
	}

	now := time.Now().UTC()
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	feed, err := s.DB.CreateFeed(context.Background(), params)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgUniqueViolation {
			return fmt.Errorf("feed url %q already exists", cmd.Args[1])
		}
		return fmt.Errorf("create feed: %w", err)
	}

	fmt.Printf("Feed added:\n")
	fmt.Printf("  name: %s\n", feed.Name)
	fmt.Printf("  url:  %s\n", feed.Url)
	fmt.Printf("  user: %s\n", user.Name)
	return nil
}
