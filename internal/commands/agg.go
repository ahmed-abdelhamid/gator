package commands

import (
	"context"
	"fmt"

	"github.com/ahmed-abdelhamid/gator/internal/cli"
	"github.com/ahmed-abdelhamid/gator/internal/rss"
)

const defaultFeedURL = "https://www.wagslane.dev/index.xml"

// Agg fetches a hard-coded RSS feed and prints its contents.
func Agg(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("agg takes no arguments")
	}

	feed, err := rss.FetchFeed(context.Background(), defaultFeedURL)
	if err != nil {
		return fmt.Errorf("fetch feed: %w", err)
	}

	fmt.Printf("Feed: %s\n", feed.Channel.Title)
	if feed.Channel.Description != "" {
		fmt.Printf("  %s\n", feed.Channel.Description)
	}
	fmt.Println()
	for _, item := range feed.Channel.Item {
		fmt.Printf("- %s\n", item.Title)
		fmt.Printf("  %s\n", item.Link)
		if item.PubDate != "" {
			fmt.Printf("  %s\n", item.PubDate)
		}
	}
	return nil
}
