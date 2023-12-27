package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (q *Queries) CountLikes(ctx context.Context, linkID uuid.UUID) (int64, error) {
	query := "SELECT COUNT(*) FROM likes WHERE link_id = $1"
	row := q.db.QueryRow(ctx, query, linkID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

type CreateLikeParams struct {
	UserID uuid.UUID
	LinkID uuid.UUID
}

func (q *Queries) CreateLike(ctx context.Context, arg CreateLikeParams) error {
	query := "INSERT INTO likes (user_id, link_id) VALUES ($1, $2)"
	_, err := q.db.Exec(ctx, query, arg.UserID, arg.LinkID)
	return err
}

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

type DeleteVoteParams struct {
	UserID uuid.UUID
	LinkID uuid.UUID
}

func (q *Queries) DeleteVote(ctx context.Context, arg DeleteVoteParams) error {
	query := "DELETE FROM likes WHERE user_id = $1 AND link_id = $2"
	_, err := q.db.Exec(ctx, query, arg.UserID, arg.LinkID)
	return err
}

type LinkFeedParams struct {
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

type LinkFeedRow struct {
	ID           uuid.UUID
	Title        string
	Url          string
	CreatedAt    pgtype.Timestamptz
	Username     string
	CommentCount int64
	LikeCount    int64
	UserLiked    int32
}

func (q *Queries) LinkFeed(ctx context.Context, arg LinkFeedParams) ([]LinkFeedRow, error) {
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
	var items []LinkFeedRow
	for rows.Next() {
		var i LinkFeedRow
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

func (q *Queries) Link(ctx context.Context, params LinkParams) (LinkFeedRow, error) {
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
	row := q.db.QueryRow(ctx, query, params.LinkID, params.UserID)
	var linkRow LinkFeedRow
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
		return LinkFeedRow{}, err
	}
	return linkRow, nil
}

type UserLikedParams struct {
	UserID uuid.UUID
	LinkID uuid.UUID
}

func (q *Queries) UserLiked(ctx context.Context, arg UserLikedParams) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND link_id = $2)"
	row := q.db.QueryRow(ctx, query, arg.UserID, arg.LinkID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

type UserFeedParams struct {
	UserID   uuid.UUID
	Username string
	Limit    int32
	Offset   int32
}

func (q *Queries) UserFeed(ctx context.Context, arg UserFeedParams) ([]LinkFeedRow, error) {
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
	var items []LinkFeedRow
	for rows.Next() {
		var i LinkFeedRow
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
