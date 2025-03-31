-- +goose Up
CREATE TABLE feed_follows (
    id integer PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    feed_id integer NOT NULL References feeds ON DELETE CASCADE,
    user_id integer NOT NULL References users ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (feed_id) REFERENCES feeds (id),
    UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
