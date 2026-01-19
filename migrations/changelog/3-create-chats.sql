--liquibase formatted sql

--changeset id:create-chats-table author:dvalkov

CREATE TABLE IF NOT EXISTS messenger.chats(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS messenger.chats_users(
    user_nickname VARCHAR REFERENCES messenger.users,
    chat_id BIGINT REFERENCES messenger.chats,
    role VARCHAR NOT NULL,
    PRIMARY KEY(user_nickname, chat_id)
);

--rollback DROP TABLE IF EXISTS messenger.chats_users;
--rollback DROP TABLE IF EXISTS messenger.chats;