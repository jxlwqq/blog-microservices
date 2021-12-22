create table posts
(
    id             bigint unsigned not null auto_increment primary key,
    user_id        bigint unsigned not null,
    title          varchar(255)    not null,
    content        text            not null,
    comments_count int unsigned    not null default 0,
    created_at     timestamp       not null default current_timestamp,
    updated_at     timestamp       not null default current_timestamp on update current_timestamp,
    index (user_id)
);