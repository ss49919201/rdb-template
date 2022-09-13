DROP TABLE IF EXISTS users;

create table
    IF not exists users (
        `id` CHAR(36) COLLATE utf8mb4_general_ci NOT NULL COMMENT '必ず36バイト確保される(足りない場合は空白で埋められる)',
        `name` VARCHAR(20) NOT NULL,
        `count` INT NOT NULL,
        `updated_at` DATETIME DEFAULT NULL,
        PRIMARY KEY (`id`),
        UNIQUE `unique_name_count` (`name`, `count`)
    ) DEFAULT CHARSET = utf8 COLLATE = utf8_bin;