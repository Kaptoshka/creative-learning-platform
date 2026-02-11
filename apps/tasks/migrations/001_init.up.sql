CREATE TABLE IF NOT EXISTS widgets(
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    type VARCHAR(255) NOT NULL,
    version INTEGER NOT NULL,

    config_schema JSONB NOT NULL DEFAULT '{}',
    submission_schema JSONB NOT NULL DEFAULT '{}',
    UNIQUE (type, version)
);

CREATE TABLE IF NOT EXISTS assignment_templates (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    creator_id UUID NOT NULL, -- Teacher ID

    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,

    widget_id UUID NOT NULL REFERENCES widgets(id),
    widget_config JSONB NOT NULL DEFAULT '{}',

    due_date TIMESTAMP NOT NULL,    -- Soft deadline
    cutoff_date TIMESTAMP NOT NULL, -- Hard deadline

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS assignment_targets (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    template_id UUID NOT NULL REFERENCES assignment_templates(id),

    group_id UUID,
    student_id UUID,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT check_target_not_empty CHECK (group_id IS NOT NULL OR student_id IS NOT NULL)
);

CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    template_id UUID NOT NULL REFERENCES assignment_templates(id),
    student_id UUID NOT NULL,

    status VARCHAR(60) NOT NULL DEFAULT 'IN_PROGRESS',
    -- 'IN_PROGRESS', 'SUBMITTED', 'GRADED', 'RETURNED'

    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    submitted_at TIMESTAMP,

    CONSTRAINT uniq_student_template UNIQUE (student_id, template_id)
);
CREATE INDEX idx_submissions_template_status ON submissions (template_id, status);
CREATE INDEX idx_submissions_student ON submissions (student_id);

CREATE TABLE IF NOT EXISTS submission_versions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    submission_id UUID NOT NULL REFERENCES submissions(id),
    version_number INTEGER NOT NULL,
    payload JSONB NOT NULL,

    time_spent_seconds INTEGER NOT NULL DEFAULT 0,

    is_autosave BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_submission_versions_latest ON submission_versions (submission_id, version_number DESC);

CREATE TABLE IF NOT EXISTS feedbacks (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    version_id UUID NOT NULL REFERENCES submission_versions(id),
    grader_id UUID NOT NULL,
    text_content TEXT,
    payload JSONB,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
