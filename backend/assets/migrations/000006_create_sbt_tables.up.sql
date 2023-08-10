CREATE TABLE IF NOT EXISTS sbt_collections (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    raw_address TEXT NOT NULL UNIQUE,
    friendly_address TEXT NOT NULL UNIQUE,
    next_item_index INTEGER NOT NULL,
    content_uri TEXT NOT NULL UNIQUE,
    raw_owner_address TEXT NOT NULL,
    friendly_owner_address TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    image TEXT,
    content_json JSONB,
    default_weight INTEGER NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    version BIGINT NOT NULL DEFAULT 1
);

CREATE INDEX sbt_collections_raw_address_idx ON sbt_collections (raw_address);

CREATE INDEX sbt_collections_friendly_address_idx ON sbt_collections (friendly_address);

CREATE TABLE IF NOT EXISTS sbt_tokens (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    raw_address TEXT NOT NULL UNIQUE,
    friendly_address TEXT NOT NULL UNIQUE,
    sbt_collections_id BIGINT NOT NULL REFERENCES sbt_collections(id),
    content_uri TEXT NOT NULL,
    raw_owner_address TEXT NOT NULL,
    friendly_owner_address TEXT NOT NULL,
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    name TEXT NOT NULL,
    description TEXT,
    image TEXT,
    content_json JSONB,
    weight INTEGER NOT NULL,
    index INTEGER NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    version BIGINT NOT NULL DEFAULT 1
);

-- add is_pinned by alter table 


CREATE INDEX sbt_tokens_raw_address_idx ON sbt_tokens (raw_address);

CREATE INDEX sbt_tokens_friendly_address_idx ON sbt_tokens (friendly_address);

CREATE INDEX sbt_tokens_collection_id_idx ON sbt_tokens (sbt_collections_id);


CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    content TEXT,
    created_at BIGINT NOT NULL
);

CREATE INDEX notifications_user_id_idx ON notifications (user_id);





