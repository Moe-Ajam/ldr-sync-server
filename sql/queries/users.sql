-- name: CreateUser :one
insert into users(id, created_at, updated_at, name, email, password)
values ($1, $2, $3, $4, $5, $6)
returning *;
