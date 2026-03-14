-- name: CreateFeed :one
INSERT INTO feed (name, url, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFeeds :many
SELECT feed.name AS feed_name, feed.url, users.name AS user_name
FROM feed
JOIN users ON feed.user_id = users.id;

-- name: GetFeedByURL :one
SELECT *
FROM feed
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feed
SET
last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :many
SELECT *
FROM feed
WHERE user_id = $1
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;
