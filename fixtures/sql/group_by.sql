SELECT user_id,count(id) FROM tasks GROUP BY user_id;

SELECT
    name,
    count(name),
    sum(count),
    max(count),
    min(count)
FROM users
GROUP BY name;