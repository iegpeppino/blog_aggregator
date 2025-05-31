-- name: CreateFeedFollow :one

WITH inserted_feed_follow as (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
) 
SELECT 
inserted_feed_follow.*,
feeds.name AS feed_name,
users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsForUser :many

SELECT 
feed_follows.*, f.name AS feed_name, u.name AS user_name
FROM feed_follows
INNER JOIN feeds f ON feed_follows.feed_id = f.id
INNER JOIN users u ON feed_follows.user_id = u.id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec

DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;