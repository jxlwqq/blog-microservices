create database if not exists comments;
use comments;
create table if not exists comments
(
    id         bigint unsigned not null auto_increment primary key,
    user_id    bigint unsigned not null,
    post_id    bigint unsigned not null,
    content    text            not null,
    created_at timestamp       not null default current_timestamp,
    updated_at timestamp       not null default current_timestamp on update current_timestamp,
    index (user_id),
    index (post_id)
);