package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateCommentParams struct {
	UserID  uuid.UUID
	LinkID  uuid.UUID
	Content string
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) error {
	query := `
		INSERT INTO comments (user_id, link_id, content)
		VALUES ($1, $2, $3)
	`
	_, err := q.db.Exec(ctx, query, arg.UserID, arg.LinkID, arg.Content)
	return err
}

type CommentRow struct {
	ID        uuid.UUID
	Username  string
	Content   string
	CreatedAt pgtype.Timestamptz
}

type CommentFeedParams struct {
	LinkID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) CommentFeed(ctx context.Context, arg CommentFeedParams) ([]CommentRow, error) {
	query := `
		SELECT 
			c.id,
			u.username,
			c.content,
			c.created_at
		FROM 
			comments c
		JOIN 
			users u ON c.user_id = u.id
		WHERE 
			c.link_id = $1
		ORDER BY 
			c.created_at DESC
		LIMIT $2
		OFFSET $3
	`
	rows, err := q.db.Query(ctx, query, arg.LinkID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commentRows []CommentRow
	for rows.Next() {
		var commentRow CommentRow
		if err := rows.Scan(
			&commentRow.ID,
			&commentRow.Username,
			&commentRow.Content,
			&commentRow.CreatedAt,
		); err != nil {
			return nil, err
		}
		commentRows = append(commentRows, commentRow)
	}

	return commentRows, nil
}
