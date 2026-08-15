CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS application_group_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id VARCHAR NOT NULL REFERENCES applications(client_id),
    group_id UUID NOT NULL REFERENCES groups(id),
    effect VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT unique_app_group UNIQUE (application_id, group_id,  effect)
);