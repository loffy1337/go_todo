-- Создание БД
create database if not exists `go_todo` character set utf8mb4 collate utf8mb4_unicode_ci;
use `go_todo`;

-- Создание таблицы пользователей
create table if not exists `users` (
    `id` int unsigned not null auto_increment,
    `username` varchar(255) not null,
    `password` varchar(255) not null,
    `avatar_url` varchar(255) null,
    `status` varchar(255) not null default "",
    `registered_at` timestamp not null default current_timestamp,
    `last_active_at` timestamp null default null,
    `is_online` tinyint(1) not null default 0,
    primary key (`id`),
    unique key `ux_users_username` (`username`)
) engine=InnoDB;

-- Создание таблицы групп задач
create table if not exists `task_groups` (
    `id` int unsigned not null auto_increment,
    `user_id` int unsigned not null,
    `title` varchar(128) not null,
    `status` enum("urgent", "daily", "longterm") not null,
    `created_at` timestamp not null default current_timestamp,
    primary key (`id`),
    key `ix_task_groups_user_id` (`user_id`),
    constraint `fk_groups_user` foreign key (`user_id`) references `users` (`id`) on delete cascade
) engine=InnoDB;

-- Создание таблицы задач
create table if not exists `tasks` (
    `id` int unsigned not null auto_increment,
    `user_id` int unsigned not null,
    `group_id` int unsigned not null,
    `title` varchar(128) not null,
    `body` text not null,
    `is_done` tinyint(1) not null default 0,
    `created_at` timestamp not null default current_timestamp,
    `updated_at` timestamp not null default current_timestamp on update current_timestamp,
    primary key (`id`),
    key `ix_tasks_user` (`user_id`),
    key `ix_tasks_group` (`group_id`),
    key `ix_tasks_done` (`is_done`),
    constraint `fk_tasks_user` foreign key (`user_id`) references `users` (`id`) on delete cascade,
    constraint `fk_tasks_group` foreign key (`group_id`) references `task_groups` (`id`) on delete cascade
) engine=InnoDB;