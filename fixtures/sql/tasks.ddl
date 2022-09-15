DROP TABLE IF EXISTS tasks;

create table
    IF not exists tasks (
        `id` char(36) COLLATE utf8mb4_general_ci NOT NULL COMMENT '必ず36バイト確保される(足りない場合は空白で埋められる)',
        `user_id` char(36) COLLATE utf8mb4_general_ci NOT NULL COMMENT '必ず36バイト確保される(足りない場合は空白で埋められる)',
        `name` VARCHAR(20) NOT NULL,
        `updated_at` DATETIME DEFAULT NULL,
        PRIMARY KEY (`id`)
    ) DEFAULT CHARSET = utf8 COLLATE = utf8_bin;