CREATE TABLE IF NOT EXISTS submission_versions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    submission_id UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    payload JSONB NOT NULL,

    time_spent_seconds INTEGER NOT NULL DEFAULT 0,

    is_autosave BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uniq_submission_version UNIQUE (submission_id, version_number)
);

CREATE INDEX idx_submission_versions_latest
    ON submission_versions (submission_id, version_number DESC);
