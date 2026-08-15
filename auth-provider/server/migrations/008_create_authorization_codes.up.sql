CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS authorization_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code_hash VARCHAR NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    application_id VARCHAR NOT NULL REFERENCES applications(client_id),
    sso_session_id UUID NOT NULL REFERENCES sso_sessions(id),
    redirect_uri TEXT NOT NULL,
    status VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP + INTERVAL '5 minutes') NOT NULL,
    used_at TIMESTAMPTZ
);