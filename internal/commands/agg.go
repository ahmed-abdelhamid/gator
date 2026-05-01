package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ahmed-abdelhamid/gator/internal/cli"
	"github.com/ahmed-abdelhamid/gator/internal/rss"
)

// Agg runs the feed scraper on a fixed interval until the process is
// signaled (Ctrl-C or SIGTERM). The argument is a Go duration like
// "1m" or "30s" understood by time.ParseDuration.
func Agg(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	interval, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("parse interval: %w", err)
	}
	if interval <= 0 {
		return fmt.Errorf("interval must be > 0; got %s", interval)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Printf("Collecting feeds every %s\n", interval)

	// Scrape once immediately, then on every tick. select-on-ctx exits
	// cleanly on Ctrl-C / SIGTERM so deferred cleanup in main() can run.
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		if err := scrapeFeeds(ctx, s); err != nil {
			log.Printf("scrape: %v", err)
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

// scrapeFeeds picks the feed with the oldest (or never-set) last_fetched_at,
// marks it fetched, downloads its RSS payload, and prints the item titles.
func scrapeFeeds(ctx context.Context, s *cli.State) error {
	feed, err := s.DB.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("get next feed: %w", err)
	}

	feed, err = s.DB.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		return fmt.Errorf("mark feed fetched: %w", err)
	}

	rssFeed, err := rss.FetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("fetch feed %s: %w", feed.Url, err)
	}

	fmt.Printf("Scraping %s (%d items)\n", feed.Name, len(rssFeed.Channel.Item))
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("- %s\n", item.Title)
	}
	return nil
}
