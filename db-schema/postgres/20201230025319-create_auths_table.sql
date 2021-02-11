
-- +migrate Up
CREATE TABLE auths (
    user_id BIGINT NOT NULL DEFAULT 0,
    type VARCHAR NOT NULL DEFAULT '',
    secret TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +migrate Down
DROP TABLE auths;
