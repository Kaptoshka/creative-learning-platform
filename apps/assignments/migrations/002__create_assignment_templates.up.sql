CREATE TABLE IF NOT EXISTS assignment_templates (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    creator_id UUID NOT NULL, -- Teacher ID

    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,

    widget_id UUID NOT NULL REFERENCES widgets(id),
    widget_config JSONB NOT NULL DEFAULT '{}',

    due_date TIMESTAMPZ NOT NULL,

    created_at TIMESTAMPZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPZ DEFAULT CURRENT_TIMESTAMP
);
