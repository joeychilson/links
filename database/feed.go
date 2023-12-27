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
	VoteScore    int64
	UserVoted    int32
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
			COALESCE(c.comment_count, 0) AS comment_count,
			COALESCE(v.vote_score, 0) AS vote_score,
			COALESCE(uv.user_voted, 0) AS user_voted
		FROM 
			links l
		JOIN 
			users u ON l.user_id = u.id
		LEFT JOIN 
			(SELECT 
				link_id, 
				COUNT(*) AS comment_count
			FROM 
				comments
			GROUP BY link_id
			) c ON l.id = c.link_id
		LEFT JOIN 
			(SELECT 
				link_id, 
				SUM(vote) AS vote_score
			FROM 
				link_votes 
			GROUP BY link_id
			) v ON l.id = v.link_id
		LEFT JOIN 
			(SELECT 
				link_id, 
				vote AS user_voted 
			FROM 
				link_votes 
				WHERE 
				user_id = $1::uuid
			) uv ON l.id = uv.link_id
		GROUP BY 
			l.id, u.username, uv.user_voted, c.comment_count, v.vote_score
		ORDER BY 
			vote_score DESC, comment_count DESC, l.created_at DESC
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
			&i.VoteScore,
			&i.UserVoted,
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
			COALESCE(comment_count, 0) AS comment_count,
			COALESCE(vote_score, 0) AS vote_score,
			COALESCE(uv.user_voted, 0) AS user_voted
		FROM 
			links l
		JOIN 
			users u ON l.user_id = u.id
		LEFT JOIN 
			(SELECT 
				link_id, 
				COUNT(*) AS comment_count
			 FROM 
				comments
			 GROUP BY link_id
			) c ON l.id = c.link_id
		LEFT JOIN 
			(SELECT 
				link_id, 
				SUM(vote) AS vote_score
			 FROM 
				link_votes 
			 GROUP BY link_id
			) v ON l.id = v.link_id
		LEFT JOIN 
			(SELECT 
				link_id, 
				vote AS user_voted 
			 FROM 
				link_votes 
			 WHERE 
				user_id = $1::uuid
			) uv ON l.id = uv.link_id
		WHERE 
			l.id = $2::uuid
		GROUP BY 
			l.id, u.username, uv.user_voted, c.comment_count, v.vote_score
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
		&linkRow.VoteScore,
		&linkRow.UserVoted,
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
			COALESCE(c.comment_count, 0) AS comment_count,
			COALESCE(v.vote_score, 0) AS vote_score,
			COALESCE(uv.user_voted, 0) AS user_voted
		FROM 
			links l
		JOIN 
			users u ON l.user_id = u.id
		LEFT JOIN 
			(SELECT 
				link_id, 
				COUNT(*) AS comment_count
			FROM 
				comments
			GROUP BY link_id
			) c ON l.id = c.link_id
		LEFT JOIN 
			(SELECT 
				link_id, 
				SUM(vote) AS vote_score
			FROM 
				link_votes 
			GROUP BY link_id
			) v ON l.id = v.link_id
		LEFT JOIN 
			(SELECT 
				link_id, 
				vote AS user_voted 
			FROM 
				link_votes 
				WHERE 
				user_id = $1::uuid
			) uv ON l.id = uv.link_id
		WHERE 
			u.username = $2
		GROUP BY 
			l.id, u.username, uv.user_voted, c.comment_count, v.vote_score
		ORDER BY 
			vote_score DESC, comment_count DESC, l.created_at DESC
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
			&i.VoteScore,
			&i.UserVoted,
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
			COALESCE(c.comment_count, 0) AS comment_count,
			COALESCE(v.vote_score, 0) AS vote_score,
			COALESCE(uv.user_voted, 0) AS user_voted
		FROM 
			links l
		JOIN 
			users u ON l.user_id = u.id
		LEFT JOIN 
			(SELECT 
				link_id, 
				COUNT(*) AS comment_count
			FROM 
				comments
			GROUP BY link_id
			) c ON l.id = c.link_id
		LEFT JOIN 
			(SELECT 
				link_id, 
				SUM(vote) AS vote_score,
				user_id
			FROM 
				link_votes 
			GROUP BY link_id, user_id
			) v ON l.id = v.link_id
		LEFT JOIN 
			(SELECT 
				link_id, 
				vote AS user_voted 
			FROM 
				link_votes 
				WHERE 
				user_id = $1::uuid
			) uv ON l.id = uv.link_id
		WHERE 
			v.user_id = $1::uuid
		GROUP BY 
			l.id, u.username, uv.user_voted, c.comment_count, v.vote_score
		ORDER BY 
			vote_score DESC, comment_count DESC, l.created_at DESC
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
			&i.VoteScore,
			&i.UserVoted,
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
