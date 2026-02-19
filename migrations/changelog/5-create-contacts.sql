--liquibase formatted sql

--changeset id:create-contacts-table author:dvalkov

CREATE TABLE IF NOT EXISTS messenger.contacts(
    user_id UUID REFERENCES messenger.users,
    contact_user_id UUID REFERENCES messenger.users,
    PRIMARY KEY(user_id, contact_user_id)
);

--rollback DROP TABLE IF EXISTS messenger.contacts;