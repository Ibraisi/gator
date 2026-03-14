-- +goose Up
CREATE TABLE feed(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name TEXT,
    url TEXT UNIQUE,
    user_id UUID ,
   CONSTRAINT fk_users
      FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE);

-- +goose Down
DROP TABLE feed;
