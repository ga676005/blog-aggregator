-- name: CreateFeedFollow :one
INSERT INTO users_feeds (id, user_id, feed_id, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteFeedFollow :exec
DELETE FROM users_feeds
WHERE id = $1;

-- name: GetUserFeedsFollow :many
SELECT
    feeds.id AS feed_id,
    feeds.created_at AS feed_created_at,
    feeds.updated_at AS feed_updated_at,
    feeds.name AS feed_name,
    feeds.url AS feed_url,
    feeds.user_id AS feed_user_id,
    users_feeds.id AS feed_follow_id,
    users_feeds.created_at AS feed_follow_at 
FROM users_feeds JOIN feeds 
ON users_feeds.feed_id = feeds.id
WHERE users_feeds.user_id = $1;