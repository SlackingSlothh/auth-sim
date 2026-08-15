CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR NOT NULL,
    actor_id UUID,
    user_id UUID REFERENCES users(id),
    application_id VARCHAR REFERENCES applications(client_id),
    session_id UUID REFERENCES sso_sessions(id),
    result VARCHAR NOT NULL,
    metadata JSONB,
    ip_address VARCHAR,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);