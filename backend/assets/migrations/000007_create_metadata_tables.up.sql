CREATE TABLE IF NOT EXISTS collection_metadata (
    id bigserial PRIMARY KEY,
    base64 text NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    image text NOT NULL,
    cover_image text,
    external_url text NOT NULL,
    marketplace text NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    version integer NOT NULL DEFAULT 1
);


-- unique
CREATE UNIQUE INDEX IF NOT EXISTS collection_metadata_base64_uindex ON collection_metadata (base64);

CREATE TABLE IF NOT EXISTS nft_metadata (
    id bigserial PRIMARY KEY,
    base64 text NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    attributes jsonb,
    external_url text NOT NULL,
    image text NOT NULL,
    marketplace text NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    version integer NOT NULL DEFAULT 1
);

-- unique

CREATE UNIQUE INDEX IF NOT EXISTS nft_metadata_base64_uindex ON nft_metadata (base64);

