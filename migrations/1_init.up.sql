CREATE TABLE IF NOT EXISTS users (
    ID integer primary key
    email text not null unique
    pass_hash BLOB not null
);

CREATE INDEX IF NOT EXISTS idx_email on users(email)

CREATE TABLE IF NOT EXISTS apps (
    ID integer primary key
    name text not null unique
    secret text not null unique
);
