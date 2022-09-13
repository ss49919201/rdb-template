SELECT
    id,
    '' as user_id,
    name,
    count,
    updated_at
FROM users
UNION
SELECT
    id,
    user_id,
    name,
    0,
    updated_at
FROM tasks;