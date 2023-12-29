package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateUserParams struct {
	Username string
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error) {
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"

	row := q.db.QueryRow(ctx, query, arg.Username, arg.Email, arg.Password)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

type CreateUserTokenParams struct {
	UserID  uuid.UUID
	Token   string
	Context string
}

func (q *Queries) CreateUserToken(ctx context.Context, arg CreateUserTokenParams) (string, error) {
	query := "INSERT INTO user_tokens (user_id, token, context) VALUES ($1, $2, $3) RETURNING token"

	row := q.db.QueryRow(ctx, query, arg.UserID, arg.Token, arg.Context)
	var token string
	err := row.Scan(&token)
	return token, err
}

type DeleteUserTokenParams struct {
	Token   string
	Context string
}

func (q *Queries) DeleteUserToken(ctx context.Context, arg DeleteUserTokenParams) error {
	query := "DELETE FROM user_tokens WHERE token = $1 AND context = $2"

	_, err := q.db.Exec(ctx, query, arg.Token, arg.Context)
	return err
}

type User struct {
	ID          uuid.UUID
	Username    string
	Email       string
	Password    string
	ConfirmedAt pgtype.Timestamptz
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

func (q *Queries) UserByEmail(ctx context.Context, email string) (User, error) {
	query := "SELECT id, username, email, password, confirmed_at, created_at, updated_at FROM users WHERE email = $1"

	row := q.db.QueryRow(ctx, query, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.ConfirmedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type UserByIDRow struct {
	ID          uuid.UUID
	Username    string
	Email       string
	ConfirmedAt pgtype.Timestamptz
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

func (q *Queries) UserByID(ctx context.Context, id uuid.UUID) (UserByIDRow, error) {
	query := "SELECT id, username, email, confirmed_at, created_at, updated_at FROM users WHERE id = $1"

	row := q.db.QueryRow(ctx, query, id)
	var i UserByIDRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.ConfirmedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type UserIDFromTokenParams struct {
	Token   string
	Context string
}

func (q *Queries) UserIDFromToken(ctx context.Context, arg UserIDFromTokenParams) (uuid.UUID, error) {
	query := "SELECT user_id FROM user_tokens WHERE token = $1 AND context = $2"

	row := q.db.QueryRow(ctx, query, arg.Token, arg.Context)
	var user_id uuid.UUID
	err := row.Scan(&user_id)
	return user_id, err
}

func (q *Queries) EmailExists(ctx context.Context, email string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"

	row := q.db.QueryRow(ctx, query, email)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

func (q *Queries) UsernameExists(ctx context.Context, username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"

	row := q.db.QueryRow(ctx, query, username)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
