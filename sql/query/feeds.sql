-- name: CreateFeed :one
INSERT INTO feeds (
  id,
  name,
  url,
  created_at,
  updated_at,
  user_id
)
VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;


-- name: GetAllFeeds :many
SELECT * FROM feeds
ORDER BY created_at DESC;

-- name: GetFeedsByUserID :many
SELECT * FROM feeds
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = NOW()
WHERE id = $1
RETURNING *;
