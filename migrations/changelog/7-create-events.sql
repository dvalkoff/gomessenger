--liquibase formatted sql

--changeset id:create-events-table author:dvalkov

CREATE TABLE IF NOT EXISTS messenger.events(
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES messenger.users,
    receiver_id UUID REFERENCES messenger.users,
    chat_id UUID REFERENCES messenger.chats,
    payload BYTEA
);

--rollback DROP TABLE IF EXISTS messenger.events;