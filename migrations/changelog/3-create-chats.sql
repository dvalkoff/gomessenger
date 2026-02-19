--liquibase formatted sql

--changeset id:create-chats-table author:dvalkov

CREATE TABLE IF NOT EXISTS messenger.chats(
    id UUID PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS messenger.chats_users(
    user_id UUID REFERENCES messenger.users,
    chat_id UUID REFERENCES messenger.chats,
    PRIMARY KEY(user_id, chat_id)
);

--rollback DROP TABLE IF EXISTS messenger.chats_users;
--rollback DROP TABLE IF EXISTS messenger.chats;