-- +goose Up
create table users (
  id text primary key,
  created_at datetime not null default current_timestamp,
  updated_at datetime not null default current_timestamp,
  name text not null unique,
  email text not null unique,
  password text not null
);

-- +goose Down
drop table users;
