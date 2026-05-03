DROP TRIGGER IF EXISTS set_updated_at ON assignment_templates;
DROP TRIGGER IF EXISTS set_updated_at ON assignment_targets;
DROP TRIGGER IF EXISTS set_updated_at ON submission_versions;
DROP TRIGGER IF EXISTS set_updated_at ON feedbacks;

DROP FUNCTION IF EXISTS trigger_set_updated_at();
