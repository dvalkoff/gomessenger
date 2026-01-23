--liquibase formatted sql

--changeset id:create-friends-table author:dvalkov

CREATE TABLE IF NOT EXISTS messenger.friends(
    nickname VARCHAR REFERENCES messenger.users,
    friends_nickname VARCHAR REFERENCES messenger.users,
    PRIMARY KEY(nickname, friends_nickname)
);

--rollback DROP TABLE IF EXISTS messenger.friends;