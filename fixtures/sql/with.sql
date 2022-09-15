WITH t as (
        SELECT *
        FROM tasks
    )
SELECT *
FROM users
    JOIN t
WHERE t.user_id = users.id;