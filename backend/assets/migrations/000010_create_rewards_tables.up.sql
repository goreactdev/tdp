CREATE TABLE IF NOT EXISTS rewards (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    sbt_token_id BIGINT NOT NULL REFERENCES sbt_tokens(id),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    version BIGINT NOT NULL DEFAULT 1,
    UNIQUE(user_id, sbt_token_id)
);

CREATE INDEX rewards_user_id_idx ON rewards (user_id);

CREATE INDEX rewards_sbt_token_id_idx ON rewards (sbt_token_id);

CREATE TABLE IF NOT EXISTS tg_messages (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    message_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    version BIGINT NOT NULL DEFAULT 1,
    UNIQUE(user_id, message_id, chat_id)
);

CREATE INDEX tg_messages_user_id_idx ON tg_messages (user_id);

CREATE INDEX tg_messages_message_id_idx ON tg_messages (message_id);

CREATE INDEX tg_messages_chat_id_idx ON tg_messages (chat_id);
