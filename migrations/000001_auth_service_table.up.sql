-- Create Enum type for user roles
create type role as enum ('admin', 'user', 'c-admin');

-- Create Users Table
create table if not exists users
(
    id         uuid                     default gen_random_uuid() primary key,
    phone      varchar(13) unique,
    email      varchar unique,
    password   varchar,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now(),
    deleted_at bigint                   default 0,
    unique (email, deleted_at)
);

-- Create User Profile Table
create table if not exists user_profile
(
    user_id         uuid references users (id) on delete cascade,
    first_name      varchar,
    last_name       varchar,
    username        varchar unique,
    nationality     varchar,
    bio             varchar, -- Fix typo
    role            role                     default 'user',
    profile_image   varchar                  default 'no images',
    followers_count int                      default 0,
    following_count int                      default 0,
    posts_count     int                      default 0,
    created_at      timestamp with time zone default now(),
    updated_at      timestamp with time zone default now(),
    is_active       bool                     default true
);

-- Create Follows Table
create table if not exists follows
(
    follower_id  uuid references users (id) on delete cascade,
    following_id uuid references users (id) on delete cascade,
    created_at   timestamp with time zone default now(),
    primary key (follower_id, following_id) -- To ensure no duplicate follows
);
