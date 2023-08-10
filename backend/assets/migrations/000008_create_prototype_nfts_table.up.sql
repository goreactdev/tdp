CREATE TABLE IF NOT EXISTS sbt_prototype (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    metadata_id BIGINT NOT NULL REFERENCES nft_metadata(id) ON DELETE CASCADE,
    weight INTEGER NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    version BIGINT NOT NULL DEFAULT 1
);


CREATE TABLE IF NOT EXISTS activities (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    token_threshold BIGINT NOT NULL,
    sbt_prototype_id BIGINT NOT NULL REFERENCES sbt_prototype(id),
    UNIQUE(name)
);

-- drop automatically generated 


CREATE INDEX activities_name_idx ON activities (name);


