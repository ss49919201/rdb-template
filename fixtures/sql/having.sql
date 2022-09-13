/* nameが重複しているグループのみ取得 */

/*
 SELECT * FROM `users` GROUP BY `name` HAVING COUNT(*) > 1;
 ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'rdb.users.id' which is not functionally dependent on columns in GROUP BY clause;
 this is incompatible with sql_mode=only_full_group_by
 */

/* Group by する name だけ指定できる */

SELECT `name` FROM `users` GROUP BY `name` HAVING COUNT(*) > 1;