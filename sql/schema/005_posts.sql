-- +goose Up
CREATE TABLE posts (
    id integer PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    published_at TIMESTAMP,
    feed_id integer NOT NULL References feeds ON DELETE CASCADE,
    FOREIGN KEY (feed_id) REFERENCES feeds (id)
);

-- +goose Down
DELETE TABLE posts;
