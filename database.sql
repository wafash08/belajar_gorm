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