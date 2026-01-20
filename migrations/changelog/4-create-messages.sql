--liquibase formatted sql

--changeset id:create-messages-table author:dvalkov

CREATE TABLE IF NOT EXISTS messenger.messages(
    id BIGSERIAL PRIMARY KEY,
    payload TEXT,
    sender VARCHAR REFERENCES messenger.users,
    chat_id BIGINT REFERENCES messenger.chats,
    sent_at TIMESTAMP WITH TIME ZONE
);

--rollback DROP TABLE IF EXISTS messenger.messages;
