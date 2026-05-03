CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    role VARCHAR(255) NOT NULL UNIQUE
);
INSERT INTO roles (id, role) VALUES
    (uuidv7(), 'admin'),
    (uuidv7(), 'student'),
    (uuidv7(), 'teacher')
ON CONFLICT DO NOTHING;
