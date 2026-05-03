CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    slug VARCHAR(64) NOT NULL UNIQUE,
    description VARCHAR(255),
    resource_group VARCHAR(64),

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
