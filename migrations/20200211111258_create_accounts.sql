CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    username TEXT,
    password_digest TEXT NOT NULL,
    password_salt TEXT NOT NULL,
    created TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS accounts_id_idx ON accounts(username);
