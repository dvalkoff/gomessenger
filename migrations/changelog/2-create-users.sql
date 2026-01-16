--liquibase formatted sql

--changeset id:create-users-table author:dvalkov
CREATE TABLE IF NOT EXISTS messenger.users(
    nickname VARCHAR PRIMARY KEY,
    name VARCHAR,
    password BYTEA NOT NULL
);

COMMENT ON TABLE messenger.users IS 'Table stores messenger users';
COMMENT ON COLUMN messenger.users.nickname IS 'Users unique nickname';
COMMENT ON COLUMN messenger.users.name IS 'Users name';
COMMENT ON COLUMN messenger.users.password IS 'Encrypted users password';

--rollback DROP TABLE IF EXISTS messenger.users;
