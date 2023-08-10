CREATE TABLE IF NOT EXISTS merch (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    amount INT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id),
    store VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
)
