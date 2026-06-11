CREATE SCHEMA chat;

CREATE TABLE chat.users(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE CHECK(char_length(username) between 3 and 100),
    password_hash VARCHAR(255) NOT NULL,
    is_online BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE chat.sessions(
    session_token VARCHAR(255) NOT NULL PRIMARY KEY,
    csrf_token VARCHAR(255) NOT NULL,
    user_id INT NOT NULL REFERENCES chat.users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,

    CHECK(expires_at >= created_at)
);

CREATE TABLE chat.friendships(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user1_id INT NOT NULL REFERENCES chat.users(id) ON DELETE CASCADE,
    user2_id INT NOT NULL REFERENCES chat.users(id) ON DELETE CASCADE,

    CHECK(user1_id < user2_id),
    UNIQUE(user1_id, user2_id)
);
CREATE INDEX idx_friendships_user1_id ON chat.friendships(user1_id);
CREATE INDEX idx_friendships_user2_id ON chat.friendships(user2_id);

CREATE TABLE chat.friendrequests(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    from_id INT NOT NULL REFERENCES chat.users(id) ON DELETE CASCADE,
    to_id INT NOT NULL REFERENCES chat.users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK(from_id != to_id),
    UNIQUE(from_id, to_id)
);
CREATE INDEX idx_friendrequests_from_id ON chat.friendrequests(from_id);
CREATE INDEX idx_friendrequests_to_id ON chat.friendrequests(to_id);

CREATE TABLE chat.messages(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    friendship_id INT NOT NULL REFERENCES chat.friendships(id) ON DELETE CASCADE,
    sender_id INT NOT NULL REFERENCES chat.users(id) ON DELETE CASCADE,
    content TEXT NOT NULL CHECK(char_length(content) > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_messages_friendship_id ON chat.messages(friendship_id);