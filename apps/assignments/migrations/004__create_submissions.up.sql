CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    template_id UUID NOT NULL REFERENCES assignment_templates(id),
    student_id UUID NOT NULL,

    status VARCHAR(60) NOT NULL DEFAULT 'IN_PROGRESS',
    -- 'IN_PROGRESS', 'SUBMITTED', 'GRADED', 'RETURNED'

    started_at TIMESTAMPZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    submitted_at TIMESTAMPZ,

    CONSTRAINT uniq_student_template UNIQUE (student_id, template_id),
    CONSTRAINT check_submissions_status
        CHECK (status IN ('IN_PROGRESS', 'SUBMITTED', 'GRADED', 'RETURNED'))
);

CREATE INDEX idx_submissions_template_status
    ON submissions (template_id, status);

CREATE INDEX idx_submissions_student
    ON submissions (student_id);
