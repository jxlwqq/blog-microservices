create table comments
(
    id         bigint unsigned not null auto_increment primary key,
    uuid       varchar(36)     not null unique,
    user_id    bigint unsigned not null,
    post_id    bigint unsigned not null,
    content    text            not null,
    created_at timestamp       not null default current_timestamp,
    updated_at timestamp       not null default current_timestamp on update current_timestamp,
    deleted_at timestamp       null,
    index (user_id),
    index (post_id)
);