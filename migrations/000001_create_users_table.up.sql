-- Create an extension to support generation of UUIDv4
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Index frequently queried columns for high-throughput searching
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_public_id ON users(public_id);