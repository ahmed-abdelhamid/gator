package commands

import (
	"context"
	"fmt"

	"github.com/ahmed-abdelhamid/gator/internal/cli"
	"github.com/ahmed-abdelhamid/gator/internal/database"
)

// AddFeed creates a feed owned by the currently logged-in user and
// auto-follows it. The two writes run in a single transaction so a
// failure can't leave a feed without its owner's follow row.
func AddFeed(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	ctx := context.Background()
	user, err := requireCurrentUser(ctx, s)
	if err != nil {
		return err
	}

	tx, err := s.Conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()
	qtx := s.DB.WithTx(tx)

	feed, err := qtx.CreateFeed(ctx, database.CreateFeedParams{
		Name:   cmd.Args[0],
		Url:    cmd.Args[1],
		UserID: user.ID,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("feed url %q already exists", cmd.Args[1])
		}
		return fmt.Errorf("create feed: %w", err)
	}

	if _, err := qtx.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return fmt.Errorf("create feed follow: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	fmt.Printf("Feed added:\n")
	fmt.Printf("  name: %s\n", feed.Name)
	fmt.Printf("  url:  %s\n", feed.Url)
	fmt.Printf("  user: %s\n", user.Name)
	return nil
}

// Feeds prints every registered feed with its owning user.
func Feeds(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("feeds takes no arguments")
	}

	feeds, err := s.DB.ListFeedsWithAuthor(context.Background())
	if err != nil {
		return fmt.Errorf("get feeds: %w", err)
	}

	for _, feed := range feeds {
		fmt.Printf("- name: %s\n", feed.Name)
		fmt.Printf("  url:  %s\n", feed.Url)
		fmt.Printf("  user: %s\n", feed.Author)
	}

	return nil
}

