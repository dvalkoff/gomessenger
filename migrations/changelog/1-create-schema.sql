--liquibase formatted sql

--changeset id:create-schema-messenger author:dvalkov
CREATE SCHEMA IF NOT EXISTS messenger;

--rollback DROP SCHEMA IF EXISTS messenger;
