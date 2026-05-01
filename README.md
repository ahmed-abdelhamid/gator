# gator

A small command-line RSS aggregator written in Go. Built as the capstone
for the [Boot.dev](https://www.boot.dev) "Build a Blog Aggregator in Go"
course.

`gator` lets you register feeds, follow them, run a periodic scraper, and
browse the latest posts — all backed by Postgres.

## Prerequisites

You need both of the following installed locally:

- **Go** 1.22 or newer — <https://go.dev/dl/>
- **PostgreSQL** 14 or newer — <https://www.postgresql.org/download/>

You'll also want **goose** to run the database migrations:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Install

```sh
go install github.com/ahmed-abdelhamid/gator@latest
```

This drops a `gator` binary into `$(go env GOBIN)` (or `$(go env GOPATH)/bin`).
Make sure that directory is on your `PATH`.

## Database setup

1. Create a database:

   ```sh
   createdb gator
   ```

2. Clone the repo to get the migration files:

   ```sh
   git clone https://github.com/ahmed-abdelhamid/gator.git
   cd gator
   ```

3. Run the migrations:

   ```sh
   goose -dir sql/schema postgres "postgres://YOUR_USER@localhost:5432/gator?sslmode=disable" up
   ```

## Configuration

`gator` reads a JSON config from `~/.gatorconfig.json`. Create it with
your database URL:

```json
{
  "db_url": "postgres://YOUR_USER@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

`current_user_name` is managed by `gator` itself once you log in or
register; you can leave it as `""` initially.

## Usage

```sh
gator <command> [args...]
```

A typical first session:

```sh
gator register alice                                    # creates a user and logs you in
gator addfeed "Hacker News" https://news.ycombinator.com/rss
gator addfeed "TechCrunch"  https://techcrunch.com/feed/
gator agg 30s                                           # scrape every 30 seconds; Ctrl-C to stop
gator browse 10                                         # show the 10 most recent posts you follow
```

### Commands

| Command | Description |
| --- | --- |
| `register <name>` | Create a new user and log in as them. |
| `login <name>` | Log in as an existing user. |
| `users` | List every registered user (current user marked). |
| `reset` | Wipe all users (cascades to feeds, follows, posts). |
| `addfeed <name> <url>` | Register a feed and auto-follow it. |
| `feeds` | List every registered feed with its owner. |
| `follow <url>` | Follow an existing feed by URL. |
| `unfollow <url>` | Unfollow a feed. |
| `following` | List feeds the current user follows. |
| `agg <duration>` | Scrape feeds on a fixed interval (e.g. `30s`, `1m`, `5m`). Long-running. |
| `browse [limit]` | Show recent posts from feeds you follow (default limit: 2). |

`addfeed`, `follow`, `unfollow`, `following`, and `browse` require an
active user — set one with `gator login` or `gator register`.

## Notes

- `gator agg` exits cleanly on `Ctrl-C` or `SIGTERM`.
- Posts are de-duplicated by URL, so re-scraping a feed is safe.
- All timestamps are stored in UTC.
