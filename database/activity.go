package database

import (
	"context"

	"github.com/google/uuid"
)

type UserFeedType string

const (
	UserFeedComments UserFeedType = "comments"
	UserFeedVoted    UserFeedType = "voted"
	UserFeedLinks    UserFeedType = "links"
)

type UserFeedLinksParams struct {
	UserID   uuid.UUID
	Username string
	Limit    int32
	Offset   int32
}

func (q *Queries) UserFeedLinks(ctx context.Context, arg UserFeedLinksParams) ([]FeedRow, error) {
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

type UserFeedVotedParams struct {
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) UserFeedVoted(ctx context.Context, arg UserFeedVotedParams) ([]FeedRow, error) {
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
