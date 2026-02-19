--liquibase formatted sql

--changeset id:create-users-table author:dvalkov

CREATE TABLE IF NOT EXISTS messenger.users(
    id UUID PRIMARY KEY,
    nickname VARCHAR NOT NULL UNIQUE,
    identity_pub_key BYTEA NOT NULL,
    signed_pub_key BYTEA NOT NULL,
    password BYTEA NOT NULL
);

--rollback DROP TABLE IF EXISTS messenger.users;
