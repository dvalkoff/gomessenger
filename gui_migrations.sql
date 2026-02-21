
CREATE TABLE IF NOT EXISTS workspace(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url VARCHAR,
    is_current BOOLEAN
);

CREATE TABLE IF NOT EXISTS auth(
    id UUID PRIMARY KEY,
    workspace_id INTEGER REFERENCES workspace,
    nickname VARCHAR,
    access_token VARCHAR,
    refresh_token VARCHAR,
    is_current BOOLEAN
);

CREATE TABLE IF NOT EXISTS keys(
    user_id UUID REFERENCES auth,

    key_type SMALLINT, -- 0 - one_time, 1 - signed_prekey, 2 - identity_key
    private_key BYTEA
);

CREATE TABLE IF NOT EXISTS contacts(
    user_id UUID REFERENCES auth,

    contact_id UUID,
    contact_nickname VARCHAR,
    PRIMARY KEY(user_id, contact_id)
);

CREATE TABLE IF NOT EXISTS chats(
    user_id UUID REFERENCES auth,

    chat_id UUID,
    PRIMARY KEY(user_id, chat_id)
);

CREATE TABLE IF NOT EXISTS chat_users(
    user_id UUID REFERENCES auth,

    chat_user_id UUID PRIMARY KEY,
    nickname VARCHAR,
    chat_id UUID
);

CREATE TABLE IF NOT EXISTS messages(
    user_id UUID REFERENCES auth,

    id BIGINT PRIMARY KEY,
    chat_id UUID,
    sender_id UUID,
    payload VARCHAR,
    created_at TIMESTAMP
);
