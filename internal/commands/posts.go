package commands

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/ahmed-abdelhamid/gator/internal/cli"
	"github.com/ahmed-abdelhamid/gator/internal/database"
	"github.com/ahmed-abdelhamid/gator/internal/rss"
	"github.com/google/uuid"
)

const defaultBrowseLimit = 2

// pubDateLayouts covers the common RSS pubDate formats. Real feeds drift
// between RFC1123Z, RFC1123, RFC822(Z), and ISO 8601; try each.
var pubDateLayouts = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	time.RFC3339,
	"2006-01-02T15:04:05Z",
	"2006-01-02 15:04:05",
}

// Browse prints the most recent posts from feeds the logged-in user
// follows. Optional argument is a positive integer limit.
func Browse(s *cli.State, cmd cli.Command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: browse [limit]")
	}

	limit := defaultBrowseLimit
	if len(cmd.Args) == 1 {
		n, err := strconv.Atoi(cmd.Args[0])
		if err != nil || n <= 0 {
			return fmt.Errorf("limit must be a positive integer; got %q", cmd.Args[0])
		}
		limit = n
	}

	posts, err := s.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("get posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("- %s\n", nullStringOr(post.Title, "(untitled)"))
		fmt.Printf("  %s\n", post.Url)
		if post.PublishedAt.Valid {
			fmt.Printf("  %s\n", post.PublishedAt.Time.Format(time.RFC1123))
		}
		if post.Description.Valid && post.Description.String != "" {
			fmt.Printf("  %s\n", post.Description.String)
		}
	}
	return nil
}

// postFromItem maps an RSS feed item into the params shape expected by
// CreatePost, normalizing nullable fields and trying multiple pubDate
// layouts before giving up.
func postFromItem(item rss.Item, feedID uuid.UUID) database.CreatePostParams {
	return database.CreatePostParams{
		Title:       nullString(item.Title),
		Url:         item.Link,
		Description: nullString(item.Description),
		PublishedAt: parsePubDate(item.PubDate),
		FeedID:      feedID,
	}
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func nullStringOr(ns sql.NullString, fallback string) string {
	if ns.Valid && ns.String != "" {
		return ns.String
	}
	return fallback
}

func parsePubDate(s string) sql.NullTime {
	if s == "" {
		return sql.NullTime{}
	}
	for _, layout := range pubDateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return sql.NullTime{Time: t, Valid: true}
		}
	}
	return sql.NullTime{}
}
