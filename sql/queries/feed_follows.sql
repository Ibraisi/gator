-- name: CreateFeedFollow :one
WITH inserted AS (
    INSERT INTO feeds_follows (user_id, feed_id)
    VALUES ($1, $2)
    RETURNING *
)
SELECT inserted.*, users.name AS user_name, feed.name AS feed_name
FROM inserted
JOIN users ON inserted.user_id = users.id
JOIN feed ON inserted.feed_id = feed.id;

-- name: GetFollowedFeedsNames :many
SELECT feed.name, feed.url
FROM feeds_follows
JOIN feed ON feeds_follows.feed_id = feed.id
WHERE feeds_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feeds_follows
USING feed
WHERE feeds_follows.feed_id = feed.id
  AND feeds_follows.user_id = $1
  AND feed.url = $2;

