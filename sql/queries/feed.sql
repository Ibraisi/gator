-- name: CreateFeed :one
INSERT INTO feed (name, url, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFeeds :many
SELECT feed.name, feed.url, users.name AS user_name
FROM feed
JOIN users ON feed.user_id = users.id;
