package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ahmed-abdelhamid/gator/internal/cli"
	"github.com/ahmed-abdelhamid/gator/internal/database"
)

// Follow records that the currently logged-in user follows the feed at the
// given URL. The feed must already exist; use AddFeed to register a new one.
func Follow(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: follow <url>")
	}

	ctx := context.Background()
	user, err := requireCurrentUser(ctx, s)
	if err != nil {
		return err
	}

	feed, err := s.DB.GetFeedByURL(ctx, cmd.Args[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("feed with url %q does not exist", cmd.Args[0])
		}
		return fmt.Errorf("get feed %q: %w", cmd.Args[0], err)
	}

	feedFollow, err := s.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("already following feed %q", cmd.Args[0])
		}
		return fmt.Errorf("create feed follow: %w", err)
	}

	fmt.Printf("Now following:\n")
	fmt.Printf("  feed: %s\n", feedFollow.FeedName)
	fmt.Printf("  user: %s\n", feedFollow.UserName)
	return nil
}

// Following lists every feed the currently logged-in user follows.
func Following(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("following takes no arguments")
	}

	ctx := context.Background()
	user, err := requireCurrentUser(ctx, s)
	if err != nil {
		return err
	}

	feedNames, err := s.DB.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("get feed follows for user %q: %w", user.Name, err)
	}

	for _, name := range feedNames {
		fmt.Printf("- %s\n", name)
	}
	return nil
}
