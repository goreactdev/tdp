CREATE TABLE IF NOT EXISTS linked_accounts (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id),
    telegram_user_id BIGINT,
    provider TEXT NOT NULL,
    avatar_url TEXT,
    login TEXT NOT NULL,
    access_token TEXT,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    version BIGINT NOT NULL DEFAULT 1,
    UNIQUE (user_id, provider)
);

CREATE INDEX linked_accounts_user_id_idx ON linked_accounts (user_id);

