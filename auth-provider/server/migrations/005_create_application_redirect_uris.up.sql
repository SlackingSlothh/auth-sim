CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS application_redirect_uris (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id VARCHAR NOT NULL REFERENCES applications(client_id),
    redirect_uri TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);