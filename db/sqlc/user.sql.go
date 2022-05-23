// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    hash_password,
    full_name
) VALUES (
    $1,$2,$3,$4
) RETURNING username, email, hash_password, full_name, changed_password_at, created_at
`

type CreateUserParams struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	HashPassword string `json:"hash_password"`
	FullName     string `json:"full_name"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.Email,
		arg.HashPassword,
		arg.FullName,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Email,
		&i.HashPassword,
		&i.FullName,
		&i.ChangedPasswordAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT username, email, hash_password, full_name, changed_password_at, created_at FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Email,
		&i.HashPassword,
		&i.FullName,
		&i.ChangedPasswordAt,
		&i.CreatedAt,
	)
	return i, err
}
