// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
insert into users(id, created_at, updated_at, name, email, password)
values (?, ?, ?, ?, ?, ?) returning id, created_at, updated_at, name, email, password
`

type CreateUserParams struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
	Password  string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Email,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Email,
		&i.Password,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
select id, created_at, updated_at, name, email, password from users where email = ?
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Email,
		&i.Password,
	)
	return i, err
}

const getUserByName = `-- name: GetUserByName :one
select id, created_at, updated_at, name, email, password from users where name = ?
`

func (q *Queries) GetUserByName(ctx context.Context, name string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByName, name)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Email,
		&i.Password,
	)
	return i, err
}
