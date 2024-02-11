-- +goose Up
ALTER TABLE posts
ADD tags TEXT[] NOT NULL DEFAULT '{}';
CREATE TABLE
    tags (
        id SERIAL PRIMARY KEY,
        tag TEXT UNIQUE NOT NULL,
        created_at TIMESTAMP NOT NULL
    );

-- +goose Down
ALTER TABLE posts DROP COLUMN tags;
DROP TABLE tags;