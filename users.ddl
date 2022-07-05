DROP TABLE IF EXISTS users;

create table
    IF not exists users (
        `id` char(36) COLLATE utf8mb4_general_ci NOT NULL,
        `name` VARCHAR(20) NOT NULL,
        `count` INT NOT NULL,
        `created_at` Datetime DEFAULT NULL,
        `updated_at` Datetime DEFAULT NULL,
        PRIMARY KEY (`id`)
    ) DEFAULT CHARSET = utf8 COLLATE = utf8_bin;