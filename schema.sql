CREATE TYPE valid_status AS ENUM ('new', 'accepted', 'rejected', 'error');
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR (30) NOT NULL,
    email VARCHAR (20) UNIQUE NOT NULL,
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
CREATE TABLE invoices (
    id BIGSERIAL PRIMARY KEY,
    service_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount >= 0),
    payment_status valid_status NOT NULL DEFAULT 'new',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX ON users (email);
CREATE INDEX ON transfers (from_user_id);
CREATE INDEX ON transfers (to_user_id);
CREATE INDEX ON transfers (from_user_id, to_user_id);
CREATE INDEX ON invoices (user_id);