INSERT INTO permissions (id, slug, description, resource_group) VALUES
    (uuidv7(), 'assignments:read',   'Read assignments',   'assignments'),
    (uuidv7(), 'assignments:write',  'Write assignments',  'assignments'),
    (uuidv7(), 'assignments:delete', 'Delete assignments', 'assignments'),
    (uuidv7(), 'submissions:create', 'Create submissions', 'submissions'),
    (uuidv7(), 'submissions:read',   'Read submissions',   'submissions'),
    (uuidv7(), 'feedback:read',      'Read feedback',      'feedback'),
    (uuidv7(), 'feedback:write',     'Write feedback',     'feedback'),
    (uuidv7(), 'widgets:read',       'Read widgets',       'widgets'),
    (uuidv7(), 'widgets:write',      'Write widgets',      'widgets')
ON CONFLICT (slug) DO NOTHING;

-- role_permissions через JOIN чтобы не хардкодить UUID
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE
    (r.role = 'admin')
    OR (r.role = 'teacher' AND p.slug IN (
        'assignments:read', 'assignments:write', 'assignments:delete',
        'submissions:read', 'feedback:read', 'feedback:write', 'widgets:read'
    ))
    OR (r.role = 'student' AND p.slug IN (
        'assignments:read', 'submissions:create',
        'submissions:read', 'feedback:read', 'widgets:read'
    ))
ON CONFLICT DO NOTHING;
