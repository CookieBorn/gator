-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;


-- name: GetUser :one
SELECT * from users
where id=$1;

-- name: GetUserId :one
SELECT id from users
where name=$1;

-- name: GetUserName :one
SELECT * from users
where name=$1;

-- name: Reset :exec
DELETE from users;
DELETE from feeds;

-- name: GetUsers :many
SELECT name from users;

-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * from feeds;
