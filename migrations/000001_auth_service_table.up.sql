CREATE TYPE role AS ENUM ('admin', 'user', 'c-admin');

CREATE TABLE IF NOT EXISTS users
(
    id         UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    email      VARCHAR UNIQUE,
    first_name VARCHAR,
    last_name  VARCHAR,
    phone      VARCHAR(13) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    deleted_at BIGINT                   DEFAULT 0
);

CREATE TABLE IF NOT EXISTS user_profile
(
    user_id       UUID REFERENCES users (id) ON DELETE CASCADE,
    username      VARCHAR UNIQUE,
    nationality   VARCHAR,
    bio           VARCHAR,
    role          role                     DEFAULT 'user',
    profile_image VARCHAR                  DEFAULT 'no image',
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT now(),
    is_active     BOOLEAN                  DEFAULT TRUE,
    UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS follows
(
    followers_id UUID REFERENCES users (id) ON DELETE CASCADE,
    following_id UUID REFERENCES users (id) ON DELETE CASCADE,
    followed_at  TIMESTAMP WITH TIME ZONE DEFAULT now(),
    PRIMARY KEY (followers_id, following_id)
);