DELETE FROM widgets
WHERE (type, version) IN (
    ('FORM', 1),
    ('TEST', 1)
);
