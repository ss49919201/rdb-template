select user_id,count(id) from tasks group by user_id;

select
    name,
    count(name),
    sum(count),
    max(count),
    min(count)
from users
group by name;