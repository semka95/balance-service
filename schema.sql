CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR (30) NOT NULL,
    email VARCHAR (20) NOT NULL,
    balance DECIMAL(10, 2) NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE TABLE transfers (
    id BIGSERIAL PRIMARY KEY,
    from_user_id BIGINT NOT NULL REFERENCES users,
    to_user_id BIGINT NOT NULL REFERENCES users,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX ON users (email);
CREATE INDEX ON transfers (from_user_id);
CREATE INDEX ON transfers (to_user_id);
CREATE INDEX ON transfers (from_user_id, to_user_id);