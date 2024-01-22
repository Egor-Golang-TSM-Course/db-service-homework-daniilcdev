-- +goose Up
CREATE TABLE post_tags (
        id SERIAL PRIMARY KEY,
        tag TEXT NOT NULL,
        post_id SERIAL NOT NULL REFERENCES posts (id) ON DELETE CASCADE
    );

-- +goose Down
DROP TABLE post_tags;