CREATE TABLE users
(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    create_at TIMESTAMPTZ NOT NULL,
    login VARCHAR(255) UNIQUE,
    password VARCHAR(80) NOT NULL,
    bill integer DEFAULT 0
);

CREATE INDEX index_idx_users ON users (id, login);