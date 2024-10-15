-- name: CreateUser :one
insert into users(id, created_at, updated_at, name, email, password)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: GetUserByName :one
select * from users where name = $1;
