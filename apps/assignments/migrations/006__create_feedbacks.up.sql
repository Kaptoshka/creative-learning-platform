CREATE TABLE IF NOT EXISTS feedbacks (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    version_id UUID NOT NULL REFERENCES submission_versions(id),
    grader_id UUID NOT NULL,
    text_content TEXT,
    payload JSONB,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPZ DEFAULT CURRENT_TIMESTAMP
);
