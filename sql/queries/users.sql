-- name: CreateUser :one
insert into users(id, created_at, updated_at, name, email, password)
values (?, ?, ?, ?, ?, ?) returning *;

-- name: GetUserByEmail :one
select * from users where email = ?;

-- name: GetUserByName :one
select * from users where name = ?;
