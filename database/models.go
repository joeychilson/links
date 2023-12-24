// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID          int32
	Username    string
	Email       string
	Password    string
	ConfirmedAt pgtype.Timestamptz
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type UserToken struct {
	ID        int32
	UserID    int32
	Token     string
	Context   string
	CreatedAt pgtype.Timestamptz
	ExpiresAt pgtype.Timestamptz
}
