-- Name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, Name, url, user_id)
VALUES
(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- Name: GetFeeds :many
SELECT * FROM feeds;

-- Name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;

-- Name: GetFeedById :one
SELECT * FROM feeds
WHERE id = $1;
