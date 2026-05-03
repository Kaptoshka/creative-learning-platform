INSERT INTO apps (name, secret, description)
VALUES (
    'test-app',
    'test-secret',
    'Test application for integration tests'
) ON CONFLICT DO NOTHING;
