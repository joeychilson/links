package database

import (
	"context"

	"github.com/google/uuid"
)

type CreateLinkParams struct {
	UserID uuid.UUID
	Title  string
	Url    string
}

func (q *Queries) CreateLink(ctx context.Context, arg CreateLinkParams) error {
	query := "INSERT INTO links (user_id, title, url) VALUES ($1, $2, $3)"
	_, err := q.db.Exec(ctx, query, arg.UserID, arg.Title, arg.Url)
	return err
}

type ToggleLikeParams struct {
	UserID uuid.UUID
	LinkID uuid.UUID
}

func (q *Queries) ToggleLike(ctx context.Context, arg ToggleLikeParams) error {
	query := `
        WITH deleted AS (
            DELETE FROM link_likes WHERE user_id = $1 AND link_id = $2 RETURNING *
        )
        INSERT INTO link_likes (user_id, link_id)
        SELECT $1, $2 WHERE NOT EXISTS (SELECT 1 FROM deleted);
    `
	_, err := q.db.Exec(ctx, query, arg.UserID, arg.LinkID)
	return err
}

func (q *Queries) CountLikes(ctx context.Context, linkID uuid.UUID) (int64, error) {
	query := "SELECT COUNT(*) FROM link_likes WHERE link_id = $1"
	row := q.db.QueryRow(ctx, query, linkID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
