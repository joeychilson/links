package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateUserParams struct {
	Avatar   string
	Email    string
	Username string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error) {
	query := "INSERT INTO users (avatar, email, username, password) VALUES ($1, $2, $3, $4) RETURNING id"
	row := q.db.QueryRow(ctx, query,
		arg.Avatar,
		arg.Email,
		arg.Username,
		arg.Password,
	)
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

func (q *Queries) EmailExists(ctx context.Context, email string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	row := q.db.QueryRow(ctx, query, email)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

type UserByIDRow struct {
	ID          uuid.UUID
	Avatar      string
	Username    string
	Email       string
	ConfirmedAt pgtype.Timestamptz
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

func (q *Queries) UserByID(ctx context.Context, id uuid.UUID) (UserByIDRow, error) {
	query := "SELECT id, avatar, username, email, confirmed_at, created_at, updated_at FROM users WHERE id = $1"
	row := q.db.QueryRow(ctx, query, id)
	var i UserByIDRow
	err := row.Scan(
		&i.ID,
		&i.Avatar,
		&i.Username,
		&i.Email,
		&i.ConfirmedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type UserIDAndPasswordByEmailRow struct {
	ID       uuid.UUID
	Password string
}

func (q *Queries) UserIDAndPasswordByEmail(ctx context.Context, email string) (UserIDAndPasswordByEmailRow, error) {
	query := "SELECT id, password FROM users WHERE email = $1"
	row := q.db.QueryRow(ctx, query, email)
	var i UserIDAndPasswordByEmailRow
	err := row.Scan(&i.ID, &i.Password)
	return i, err
}

type UserIDByTokenParams struct {
	Token   string
	Context string
}

func (q *Queries) UserIDByToken(ctx context.Context, arg UserIDByTokenParams) (uuid.UUID, error) {
	query := "SELECT user_id FROM user_tokens WHERE token = $1 AND context = $2"
	row := q.db.QueryRow(ctx, query, arg.Token, arg.Context)
	var user_id uuid.UUID
	err := row.Scan(&user_id)
	return user_id, err
}

func (q *Queries) UsernameExists(ctx context.Context, username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"
	row := q.db.QueryRow(ctx, query, username)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

type UserRow struct {
	ID uuid.UUID
}

func (q *Queries) UserList(ctx context.Context) ([]UserRow, error) {
	query := "SELECT id FROM users"
	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserRow{}
	for rows.Next() {
		var i UserRow
		if err := rows.Scan(&i.ID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

type UserProfile struct {
	ID       uuid.UUID
	Avatar   string
	Username string
	Likes    int64
	Links    int64
	Comments int64
}

func (q *Queries) UserProfile(ctx context.Context, username string) (UserProfile, error) {
	query := `
		SELECT
			u.id,
			u.avatar,
			u.username,
			COUNT(DISTINCT l.id) AS likes,
			COUNT(DISTINCT lk.id) AS links,
			COUNT(DISTINCT c.id) AS comments
		FROM users u
		LEFT JOIN link_likes l ON l.user_id = u.id
		LEFT JOIN links lk ON lk.user_id = u.id
		LEFT JOIN comments c ON c.user_id = u.id
		WHERE u.username = $1
		GROUP BY u.id
	`
	row := q.db.QueryRow(ctx, query, username)
	var i UserProfile
	err := row.Scan(
		&i.ID,
		&i.Avatar,
		&i.Username,
		&i.Likes,
		&i.Links,
		&i.Comments,
	)
	return i, err
}
