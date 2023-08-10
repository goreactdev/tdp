CREATE TABLE IF NOT EXISTS stored_rewards (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_address TEXT NOT NULL REFERENCES users(friendly_address) ON DELETE CASCADE,
    collection_address TEXT NOT NULL REFERENCES sbt_collections(friendly_address) ON DELETE CASCADE,
    base64_metadata TEXT NOT NULL REFERENCES nft_metadata(base64) ON DELETE CASCADE,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    approved_by_user BOOLEAN NOT NULL DEFAULT FALSE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);


CREATE INDEX IF NOT EXISTS stored_rewards_user_address_idx ON stored_rewards(user_address);

CREATE INDEX IF NOT EXISTS stored_rewards_collection_address_idx ON stored_rewards(collection_address);

CREATE INDEX IF NOT EXISTS stored_rewards_base64_metadata_idx ON stored_rewards(base64_metadata);


