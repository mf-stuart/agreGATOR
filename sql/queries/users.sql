-- Name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, Name)
VALUES
(
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- Name: GetUser :one
SELECT * FROM users
WHERE Name = $1;

-- Name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- Name: GetUsers :many
SELECT * FROM users;

-- Name: Reset :exec
DELETE FROM users;
