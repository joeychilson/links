package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type FeedRow struct {
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

func (q *Queries) LinkFeed(ctx context.Context, arg LinkFeedParams) ([]FeedRow, error) {
	query := `
		SELECT 
			l.id AS id,
			l.title,
			l.url,
			l.created_at,
			u.username,
			COUNT(DISTINCT c.id) AS comment_count,
			COUNT(DISTINCT lk.id) AS like_count,
			COALESCE(ul.user_liked, 0) AS user_liked
		FROM 
			links l
		JOIN 
			users u ON l.user_id = u.id
		LEFT JOIN 
			comments c ON l.id = c.link_id
		LEFT JOIN 
			link_likes lk ON l.id = lk.link_id
		LEFT JOIN 
			(SELECT 
				link_id, 
				COUNT(*) AS user_liked 
			FROM 
				link_likes 
			WHERE 
				user_id = $1::uuid
			GROUP BY 
				link_id) ul ON l.id = ul.link_id
		GROUP BY 
			l.id, u.username, ul.user_liked
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
	var items []FeedRow
	for rows.Next() {
		var i FeedRow
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

func (q *Queries) Link(ctx context.Context, params LinkParams) (FeedRow, error) {
	query := `
		SELECT 
			l.id AS id,
			l.title,
			l.url,
			l.created_at,
			u.username,
			COUNT(DISTINCT c.id) AS comment_count,
			COUNT(DISTINCT lk.id) AS like_count,
			COALESCE(ul.user_liked, 0) AS user_liked
		FROM 
			links l
		JOIN 
			users u ON l.user_id = u.id
		LEFT JOIN 
			comments c ON l.id = c.link_id
		LEFT JOIN 
			link_likes lk ON l.id = lk.link_id
		LEFT JOIN 
			(SELECT 
				link_id, 
				COUNT(*) AS user_liked 
			FROM 
				link_likes 
			WHERE 
				user_id = $1::uuid
			GROUP BY 
				link_id) ul ON l.id = ul.link_id
		WHERE 
			l.id = $2::uuid
		GROUP BY 
			l.id, u.username, ul.user_liked
	`
	row := q.db.QueryRow(ctx, query, params.UserID, params.LinkID)
	var linkRow FeedRow
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
		return FeedRow{}, err
	}
	return linkRow, nil
}

type UserFeedParams struct {
	UserID   uuid.UUID
	Username string
	Limit    int32
	Offset   int32
}

func (q *Queries) UserFeed(ctx context.Context, arg UserFeedParams) ([]FeedRow, error) {
	query := `
		SELECT 
			l.id AS id,
			l.title,
			l.url,
			l.created_at,
			u.username,
			COUNT(DISTINCT c.id) AS comment_count,
			COUNT(DISTINCT lk.id) AS like_count,
			COALESCE(ul.user_liked, 0) AS user_liked
		FROM 
			links l
		JOIN 
			users u ON l.user_id = u.id
		LEFT JOIN 
			comments c ON l.id = c.link_id
		LEFT JOIN 
			link_likes lk ON l.id = lk.link_id
		LEFT JOIN 
			(SELECT 
				link_id, 
				COUNT(*) AS user_liked 
			FROM 
				link_likes 
			WHERE 
				user_id = $1::uuid
			GROUP BY 
				link_id) ul ON l.id = ul.link_id
		WHERE 
			u.username = $2
		GROUP BY 
			l.id, u.username, ul.user_liked
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
	var items []FeedRow
	for rows.Next() {
		var i FeedRow
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

func (q *Queries) LikedFeed(ctx context.Context, arg LikedFeedParams) ([]FeedRow, error) {
	query := `
		SELECT 
			l.id AS id,
			l.title,
			l.url,
			l.created_at,
			u.username,
			COUNT(DISTINCT c.id) AS comment_count,
			COUNT(DISTINCT lk.id) AS like_count,
			1 AS user_liked -- Since the query is filtered by lk.user_id = $1::uuid
		FROM 
			links l
		JOIN 
			users u ON l.user_id = u.id
		LEFT JOIN 
			comments c ON l.id = c.link_id
		JOIN 
			link_likes lk ON l.id = lk.link_id
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
	var items []FeedRow
	for rows.Next() {
		var i FeedRow
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
