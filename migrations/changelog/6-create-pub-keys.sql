--liquibase formatted sql

--changeset id:create-pub-keys-table author:dvalkov

CREATE TABLE IF NOT EXISTS messenger.users_pub_keys(
    id BIGSERIAL PRIMARY KEY,
    pub_key BYTEA NOT NULL,
    key_type SMALLINT NOT NULL,
    user_id UUID REFERENCES messenger.users
);

COMMENT ON COLUMN messenger.users_pub_keys.key_type IS '0 - Identity key; 1 - Signed prekey; 2 - one-time prekey';

--rollback DROP TABLE IF EXISTS messenger.users_pub_keys;
