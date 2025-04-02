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

-- name: GetFeedName :one
SELECT * from feeds
Where url=$1;

-- name: MarkedFetchFeed :exec
UPDATE feeds
set updated_at=$1 and last_fetched_at=$1
Where id=$2;

-- name: GetNextFeedToFetch :one
SELECT * from feeds
Order by last_fetched_at NULLS FIRST
limit 1;
