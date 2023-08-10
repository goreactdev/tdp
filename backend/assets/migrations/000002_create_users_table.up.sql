CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    username citext NOT NULL UNIQUE,
    raw_address TEXT NOT NULL UNIQUE,
    friendly_address TEXT NOT NULL UNIQUE,
    job TEXT, 
    bio TEXT,
    languages TEXT[],
    certifications TEXT[],
    avatar_url TEXT,
    awards_count INTEGER NOT NULL DEFAULT 0,
    messages_count INTEGER NOT NULL DEFAULT 0,
    rating INTEGER NOT NULL DEFAULT 0,
    last_award_at BIGINT,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    version BIGINT NOT NULL DEFAULT 1
);

CREATE INDEX users_username_idx ON users (username);

CREATE INDEX users_first_name_idx ON users (first_name);

CREATE INDEX users_last_name_idx ON users (last_name);

CREATE INDEX users_address_idx ON users (raw_address);

CREATE INDEX users_friendly_address_idx ON users (friendly_address);



CREATE SEQUENCE IF NOT EXISTS username_seq;

CREATE OR REPLACE FUNCTION generate_unique_username()
RETURNS TEXT AS $$
DECLARE
    username TEXT;
    unique_num INTEGER;
BEGIN
    -- Get a unique number from the sequence
    SELECT nextval('username_seq') INTO unique_num;

    -- Generate a random string
    SELECT md5(random()::text || clock_timestamp()::text) INTO username;
    
    -- Concatenate the random string with the unique number
    RETURN 'user_' || unique_num || '_' || substring(username, 1, 8);
END;
$$ LANGUAGE plpgsql;
