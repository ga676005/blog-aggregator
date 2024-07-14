-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (url) DO NOTHING;


-- name: GetPostsByUser :many
SELECT * FROM posts 
WHERE feed_id = ANY(sqlc.arg(feed_ids)::uuid[])
ORDER BY published_at DESC
LIMIT $1;
