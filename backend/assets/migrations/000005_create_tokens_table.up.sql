CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users on DELETE CASCADE,
    expiry bigint NOT NULL,
    scope text NOT NULL
)
