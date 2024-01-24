-- +goose Up
CREATE TABLE post_comments (
        id SERIAL PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        comment_text TEXT NOT NULL,
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        post_id SERIAL NOT NULL REFERENCES posts (id) ON DELETE CASCADE
    );

-- +goose Down
DROP TABLE post_comments;