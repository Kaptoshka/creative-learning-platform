CREATE TABLE IF NOT EXISTS assignment_targets (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    template_id UUID NOT NULL REFERENCES assignment_templates(id),

    group_id UUID,
    student_id UUID,

    created_at TIMESTAMPZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT check_target_not_empty CHECK (group_id IS NOT NULL OR student_id IS NOT NULL)
);
