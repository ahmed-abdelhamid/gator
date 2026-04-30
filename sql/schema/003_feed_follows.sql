-- +goose Up
CREATE TABLE feed_follows (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  UNIQUE (user_id, feed_id)
);

CREATE INDEX feed_follows_feed_id_idx ON feed_follows(feed_id);

-- +goose Down
DROP TABLE feed_follows;