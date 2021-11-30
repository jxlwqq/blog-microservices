create database if not exists posts;
use posts;
create table if not exists posts
(
    id         bigint unsigned not null auto_increment primary key,
    user_id    bigint unsigned not null,
    title      varchar(255)    not null,
    content    text            not null,
    created_at timestamp       not null default current_timestamp,
    updated_at timestamp       not null default current_timestamp on update current_timestamp,
    index (user_id)
    )