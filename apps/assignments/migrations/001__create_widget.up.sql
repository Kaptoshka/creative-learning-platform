CREATE TABLE IF NOT EXISTS widgets(
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    type VARCHAR(255) NOT NULL,
    version INTEGER NOT NULL,

    config_schema JSONB NOT NULL DEFAULT '{}',
    submission_schema JSONB NOT NULL DEFAULT '{}',
    UNIQUE (type, version)
);
