CREATE TABLE IF NOT EXISTS assignments (
    id UUIDv7 PRIMARY KEY,
    creator_id UUIDv7 NOT NULL,
    student_id UUIDv7 NOT NULL,
    title VARCHAR(255) NOT NULL,
    widget_id INTEGER NOT NULL REFERENCES widgets(id),
    widget_config JSONB NOT NULL,
    due_date TIMESTAMP NOT NULL,
    cutoff_date TIMESTAMP NOT NULL,
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
    id UUIDv7 PRIMARY KEY,
    assignment_id UUIDv7 NOT NULL REFERENCES assignments(id),
    creator_id UUIDv7 NOT NULL,
    status VARCHAR(60) NOT NULL,
    current_version_id UUIDv7 NOT NULL REFERENCES submission_versions(id),
    started_at TIMESTAMP NOT NULL,
    submitted_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS submission_versions (
    id UUIDv7 PRIMARY KEY,
    submission_id UUIDv7 NOT NULL REFERENCES submissions(id),
    version_number INTEGER NOT NULL,
    payload JSONB NOT NULL,
    is_late BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS feedbacks (
    id UUIDv7 PRIMARY KEY,
    submission_version_id UUIDv7 NOT NULL REFERENCES submission_attempts(id),
    grader_id UUIDv7 NOT NULL,
    feedback VARCHAR(255) NOT NULL,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
