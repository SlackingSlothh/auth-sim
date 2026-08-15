CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS access_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash VARCHAR NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    application_id VARCHAR NOT NULL REFERENCES applications(client_id),
    sso_session_id UUID NOT NULL REFERENCES sso_sessions(id),
    scopes JSONB,
    status VARCHAR NOT NULL,
    issued_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP + INTERVAL '14 days') NOT NULL,
    revoked_at TIMESTAMPTZ
);