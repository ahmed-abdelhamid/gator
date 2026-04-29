-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: ListFeedsWithAuthor :many
SELECT
  feeds.name,
  feeds.url,
  users.name AS author
FROM feeds
INNER JOIN users
ON users.id = feeds.user_id
ORDER BY feeds.created_at DESC;
