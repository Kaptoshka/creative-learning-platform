CREATE TABLE IF NOT EXISTS assignment_templates (
    id UUID PRIMARY KEY,
    creator_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    widget_id INTEGER NOT NULL REFERENCES widgets(id),
    widget_config JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS student_assignments (
    id UUID PRIMARY KEY,
    template_id UUID NOT NULL REFERENCES assignment_templates(id),
    student_id UUID NOT NULL,
    due_date TIMESTAMP NOT NULL,
    cutoff_date TIMESTAMP NOT NULL,
    status VARCHAR(60) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS widgets(
    id SERIAL PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    version INTEGER NOT NULL,
    config_schema JSONB NOT NULL,
    submission_schema JSONB NOT NULL
);

CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY,
    assignment_id UUID NOT NULL REFERENCES student_assignments(id),
    creator_id UUID NOT NULL,
    status VARCHAR(60) NOT NULL,
    current_version_id UUID REFERENCES submission_versions(id),
    started_at TIMESTAMP NOT NULL,
    submitted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS submission_versions (
    id UUID PRIMARY KEY,
    submission_id UUID NOT NULL REFERENCES submissions(id),
    version_number INTEGER NOT NULL,
    payload JSONB NOT NULL,
    is_late BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS feedbacks (
    id UUID PRIMARY KEY,
    submission_version_id UUID NOT NULL REFERENCES submission_versions(id),
    grader_id UUID NOT NULL,
    feedback TEXT NOT NULL,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
