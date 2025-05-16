-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT
    inserted_feed_follow.id,
    inserted_feed_follow.created_at,
    inserted_feed_follow.updated_at,
    inserted_feed_follow.user_id,
    inserted_feed_follow.feed_id,
    users.name AS user_name,
    feeds.name AS feed_name
FROM inserted_feed_follow
JOIN users ON users.id = inserted_feed_follow.user_id
JOIN feeds ON feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT
    feed_follows.id,
    feed_follows.created_at,
    feed_follows.updated_at,
    feed_follows.user_id,
    feed_follows.feed_id,
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follows
JOIN feeds ON feeds.id = feed_follows.feed_id
JOIN users ON users.id = feed_follows.user_id
WHERE users.name = $1;

-- name: DeleteFeedFollowByUserAndURL :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1
  AND feed_follows.feed_id = (
    SELECT feeds.id FROM feeds WHERE feeds.url = $2
  );

