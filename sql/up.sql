CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "users" (
    id uuid NOT NULL DEFAULT uuid_generate_v4() CONSTRAINT subscribers_pkey PRIMARY KEY,
    first_name varchar(255) NOT NULL,
    last_name varchar(255) NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS "users_email_key" ON "users" (email);

CREATE TABLE IF NOT EXISTS "groups" (
    id uuid NOT NULL DEFAULT uuid_generate_v4() CONSTRAINT groups_pkey PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS "groups_name_key" ON "groups" (name);

CREATE TABLE IF NOT EXISTS "users_groups" (
    user_id uuid references users NOT NULL,
    group_id uuid references groups NOT NULL,
    CONSTRAINT users_groups_pkey PRIMARY KEY (user_id, group_id)
);

CREATE TABLE IF NOT EXISTS "mails" (
    id uuid NOT NULL DEFAULT uuid_generate_v4() CONSTRAINT mails_pkey PRIMARY KEY,
    to_user_id uuid references users NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMP,
    watched BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS "mails_to_user_id_index" ON "mails" (to_user_id);
