CREATE TYPE role AS ENUM ('admin', 'user', 'c-admin');

CREATE TABLE IF NOT EXISTS users
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR UNIQUE,
    phone         VARCHAR(16) UNIQUE,
    first_name    VARCHAR NOT NULL,
    last_name     VARCHAR NOT NULL,
    username      VARCHAR,
    country       VARCHAR,
    password      VARCHAR NOT NULL,
    bio           VARCHAR,
    role          role             DEFAULT 'user',
    profile_image VARCHAR          DEFAULT 'no image',
    created_at    TIMESTAMP        DEFAULT now(),
    updated_at    TIMESTAMP        DEFAULT now(),
    deleted_at    BIGINT           DEFAULT 0,
    UNIQUE (username, deleted_at)
);

CREATE TABLE IF NOT EXISTS follows
(
    follower_id  UUID REFERENCES users (id) ON DELETE CASCADE,
    following_id UUID REFERENCES users (id) ON DELETE CASCADE,
    followed_at  TIMESTAMP DEFAULT now(),
    PRIMARY KEY (follower_id, following_id)
);