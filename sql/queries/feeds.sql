-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE id = $1;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY 
    last_fetched_at IS NULL DESC,
    last_fetched_at
LIMIT $1;

-- name: MarkFeedFetched :exec
UPDATE feeds 
SET last_fetched_at = $1 and updated_at = $1
WHERE id = $2;