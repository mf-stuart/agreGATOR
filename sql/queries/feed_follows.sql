-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
    VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.Name AS feed_name,
    users.Name AS user_name
FROM inserted_feed_follow
INNER JOIN users ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT * FROM feed_follows
WHERE user_id = $1;

-- name: DeleteFeedFollow :exec
WITH feed_id_match AS (
    SELECT id FROM feeds
    WHERE url = $1
)
DELETE FROM feed_follows
WHERE feed_follows.user_id = $2
AND feed_id = (SELECT id FROM feed_id_match);

