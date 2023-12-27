package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

type LinkRow struct {
	ID           uuid.UUID
	Title        string
	Url          string
	CreatedAt    pgtype.Timestamptz
	Username     string
	CommentCount int64
	LikeCount    int64
	UserLiked    int32
}

type LinkFeedParams struct {
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) LinkFeed(ctx context.Context, arg LinkFeedParams) ([]LinkRow, error) {
	query := `
		SELECT 
			l.id AS id,
			l.title,
			l.url,
			l.created_at,
			u.username,
			COUNT(DISTINCT c.id) AS comment_count,
			COUNT(DISTINCT lk.id) AS like_count,
			CASE 
				WHEN $1::uuid IS NOT NULL THEN SUM(CASE WHEN lk.user_id = $1::uuid THEN 1 ELSE 0 END)
				ELSE 0 
			END AS user_liked
		FROM 
			links l
		JOIN 
			users u ON l.user_id = u.id
		LEFT JOIN 
			comments c ON l.id = c.link_id
		LEFT JOIN 
			likes lk ON l.id = lk.link_id
		GROUP BY 
			l.id, u.username
		ORDER BY 
			like_count DESC, comment_count DESC, l.created_at DESC
		LIMIT 
			$2
		OFFSET 
			$3
	`
	rows, err := q.db.Query(ctx, query, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LinkRow
	for rows.Next() {
		var i LinkRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Url,
			&i.CreatedAt,
			&i.Username,
			&i.CommentCount,
			&i.LikeCount,
			&i.UserLiked,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type LinkParams struct {
	UserID uuid.UUID
	LinkID uuid.UUID
}

func (q *Queries) Link(ctx context.Context, params LinkParams) (LinkRow, error) {
	query := `
		SELECT 
			l.id AS id,
			l.title,
			l.url,
			l.created_at,
			u.username,
			COUNT(DISTINCT c.id) AS comment_count,
			COUNT(DISTINCT lk.id) AS like_count,
			CASE 
				WHEN $1::uuid IS NOT NULL THEN SUM(CASE WHEN lk.user_id = $1::uuid THEN 1 ELSE 0 END)
				ELSE 0 
			END AS user_liked
		FROM 
			links l
		JOIN 
			users u ON l.user_id = u.id
		LEFT JOIN 
			comments c ON l.id = c.link_id
		LEFT JOIN 
			likes lk ON l.id = lk.link_id
		WHERE 
			l.id = $2::uuid
		GROUP BY 
			l.id, u.username
	`
	row := q.db.QueryRow(ctx, query, params.UserID, params.LinkID)
	var linkRow LinkRow
	if err := row.Scan(
		&linkRow.ID,
		&linkRow.Title,
		&linkRow.Url,
		&linkRow.CreatedAt,
		&linkRow.Username,
		&linkRow.CommentCount,
		&linkRow.LikeCount,
		&linkRow.UserLiked,
	); err != nil {
		return LinkRow{}, err
	}
	return linkRow, nil
}

type UserFeedParams struct {
	UserID   uuid.UUID
	Username string
	Limit    int32
	Offset   int32
}

func (q *Queries) UserFeed(ctx context.Context, arg UserFeedParams) ([]LinkRow, error) {
	query := `
		SELECT 
			l.id AS id,
			l.title,
			l.url,
			l.created_at,
			u.username,
			COUNT(DISTINCT c.id) AS comment_count,
			COUNT(DISTINCT lk.id) AS like_count,
			CASE 
				WHEN $1::uuid IS NOT NULL THEN SUM(CASE WHEN lk.user_id = $1::uuid THEN 1 ELSE 0 END)
				ELSE 0 
			END AS user_liked
		FROM 
			links l
		JOIN 
			users u ON l.user_id = u.id
		LEFT JOIN 
			comments c ON l.id = c.link_id
		LEFT JOIN 
			likes lk ON l.id = lk.link_id
		WHERE 
			u.username = $2
		GROUP BY 
			l.id, u.username
		ORDER BY 
			like_count DESC, comment_count DESC, l.created_at DESC
		LIMIT 
			$3
		OFFSET 
			$4
	`
	rows, err := q.db.Query(ctx, query, arg.UserID, arg.Username, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LinkRow
	for rows.Next() {
		var i LinkRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Url,
			&i.CreatedAt,
			&i.Username,
			&i.CommentCount,
			&i.LikeCount,
			&i.UserLiked,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type LikedFeedParams struct {
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) LikedFeed(ctx context.Context, arg LikedFeedParams) ([]LinkRow, error) {
	query := `
		SELECT 
			l.id AS id,
			l.title,
			l.url,
			l.created_at,
			u.username,
			COUNT(DISTINCT c.id) AS comment_count,
			COUNT(DISTINCT lk.id) AS like_count,
			CASE 
				WHEN $1::uuid IS NOT NULL THEN SUM(CASE WHEN lk.user_id = $1::uuid THEN 1 ELSE 0 END)
				ELSE 0 
			END AS user_liked
		FROM 
			links l
		JOIN 
			users u ON l.user_id = u.id
		LEFT JOIN 
			comments c ON l.id = c.link_id
		JOIN 
			likes lk ON l.id = lk.link_id
		WHERE 
			lk.user_id = $1::uuid
		GROUP BY 
			l.id, u.username
		ORDER BY 
			like_count DESC, comment_count DESC, l.created_at DESC
		LIMIT 
			$2
		OFFSET 
			$3
	`
	rows, err := q.db.Query(ctx, query, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LinkRow
	for rows.Next() {
		var i LinkRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Url,
			&i.CreatedAt,
			&i.Username,
			&i.CommentCount,
			&i.LikeCount,
			&i.UserLiked,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type ToggleLikeParams struct {
	UserID uuid.UUID
	LinkID uuid.UUID
}

func (q *Queries) ToggleLike(ctx context.Context, arg ToggleLikeParams) error {
	query := `
        WITH deleted AS (
            DELETE FROM likes WHERE user_id = $1 AND link_id = $2 RETURNING *
        )
        INSERT INTO likes (user_id, link_id)
        SELECT $1, $2 WHERE NOT EXISTS (SELECT 1 FROM deleted);
    `
	_, err := q.db.Exec(ctx, query, arg.UserID, arg.LinkID)
	return err
}

func (q *Queries) CountLikes(ctx context.Context, linkID uuid.UUID) (int64, error) {
	query := "SELECT COUNT(*) FROM likes WHERE link_id = $1"
	row := q.db.QueryRow(ctx, query, linkID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
