CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR (30) NOT NULL,
    email VARCHAR (20) NOT NULL,
    balance DECIMAL(10, 2) NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX ON users (email);