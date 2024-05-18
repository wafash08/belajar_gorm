create table sample (
  id varchar(100) not null,
  name varchar(100) not null,
  primary key (id)
) engine = inndodb

create table users
(
  id         varchar(100) not null,
  password   varchar(100) not null,
  name       varchar(100) not null,
  created_at timestamp    not null default current_timestamp,
  updated_at timestamp    not null default current_timestamp on update current_timestamp,
  primary key (id)
) engine = InnoDB;

DELETE FROM users WHERE id = '1';

create table user_logs
(
    id         int auto_increment,
    user_id    varchar(100) not null,
    action     varchar(100) not null,
    created_at timestamp    not null default current_timestamp,
    updated_at timestamp    not null default current_timestamp on update current_timestamp,
    primary key (id)
) engine = innodb;

create table todos
(
    id          bigint       not null auto_increment,
    user_id     varchar(100) not null,
    title       varchar(100) not null,
    description text         null,
    created_at  timestamp    not null default current_timestamp,
    updated_at  timestamp    not null default current_timestamp on update current_timestamp,
    deleted_at  timestamp    null,
    primary key (id)
) engine = innodb;