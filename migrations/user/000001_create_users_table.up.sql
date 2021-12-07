create table users
(
    id         bigint unsigned not null auto_increment primary key,
    username   varchar(255)    not null unique,
    email      varchar(255)    not null unique,
    avatar     varchar(255) default null,
    password   varchar(255)    not null,
    created_at timestamp    default current_timestamp,
    updated_at timestamp    default current_timestamp on update current_timestamp
);