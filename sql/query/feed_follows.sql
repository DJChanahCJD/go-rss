-- name: CreateFeedFollow :one
INSERT INTO feed_follows (
  id,
  created_at,
  updated_at,
  user_id,
  feed_id
)
VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;

-- name: GetAllFeedFollows :many
SELECT * FROM feed_follows
ORDER BY created_at DESC;

-- name: GetFeedFollowsByUserID :many
SELECT ff.*, feeds.name as feed_name, feeds.url as feed_url FROM feed_follows ff
JOIN feeds ON ff.feed_id = feeds.id
WHERE ff.user_id = $1
ORDER BY created_at DESC;
