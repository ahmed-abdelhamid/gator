-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: ListFeedsWithAuthor :many
SELECT
  feeds.name,
  feeds.url,
  users.name AS author
FROM feeds
INNER JOIN users
ON users.id = feeds.user_id
ORDER BY feeds.created_at DESC;
