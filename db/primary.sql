CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    password_hash VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS user_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    group_id UUID NOT NULL REFERENCES groups(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT user_group_unique UNIQUE (user_id, group_id)
);

CREATE TABLE IF NOT EXISTS applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    client_id VARCHAR UNIQUE NOT NULL,
    client_secret_hash VARCHAR,
    status VARCHAR NOT NULL,
    launch_url TEXT,
    logout_notification_url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS application_redirect_uris (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id VARCHAR NOT NULL REFERENCES applications(client_id),
    redirect_uri TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS application_group_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id VARCHAR NOT NULL REFERENCES applications(client_id),
    group_id UUID NOT NULL REFERENCES groups(id),
    effect VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT unique_app_group UNIQUE (application_id, group_id,  effect)
);

CREATE TABLE IF NOT EXISTS sso_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    session_token_hash VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP + INTERVAL '14 days') NOT NULL,
    last_activity_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    revoke_reason VARCHAR,
    ip_address VARCHAR,
    user_agent TEXT
);

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