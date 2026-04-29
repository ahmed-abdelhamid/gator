// Package rss fetches and parses RSS 2.0 feeds.
package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

const maxFeedBytes = 10 << 20

var client = &http.Client{Timeout: 30 * time.Second}

// Feed is a parsed RSS 2.0 feed.
type Feed struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Item        []Item `xml:"item"`
	} `xml:"channel"`
}

// Item is a single entry within a Feed's channel.
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// FetchFeed downloads feedURL and returns the parsed feed with HTML
// entities decoded in the title and description fields.
func FetchFeed(ctx context.Context, feedURL string) (*Feed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request for %s: %w", feedURL, err)
	}
	req.Header.Set("User-Agent", "gator/0.1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http GET %s: %w", feedURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("fetch %s: status %s", feedURL, resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxFeedBytes))
	if err != nil {
		return nil, fmt.Errorf("read body from %s: %w", feedURL, err)
	}

	var feed Feed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("unmarshal feed from %s: %w", feedURL, err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}
