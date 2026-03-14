-- +goose Up
CREATE TABLE feeds_follows(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    UNIQUE(user_id, feed_id),
    CONSTRAINT fk_ff_users
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_ff_feed
        FOREIGN KEY(feed_id)
            REFERENCES feed(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds_follows;
