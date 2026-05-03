CREATE TABLE IF NOT EXISTS apps (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(255) NOT NULL UNIQUE,
    secret VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
INSERT INTO apps (id, name, secret, description) VALUES
    (uuidv7(), 'assignments-service', '', 'Assignment lifecycle microservice')
ON CONFLICT DO NOTHING;
