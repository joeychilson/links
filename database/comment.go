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
	ID         uuid.UUID
	Username   string
	Content    string
	ReplyCount int64
	VoteScore  int64
	UserVoted  int32
	CreatedAt  pgtype.Timestamptz
}

type CommentFeedParams struct {
	UserID uuid.UUID
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
            (
                SELECT COUNT(*)
                FROM comments as rc
                WHERE rc.parent_id = c.id
            ) as reply_count,
            COALESCE(SUM(cv.vote), 0) as vote_score,
            COALESCE(
                (
                    SELECT cv.vote
                    FROM comment_votes as cv
                    WHERE cv.comment_id = c.id AND cv.user_id = $4
                ),
                0
            ) as user_voted,
            c.created_at
        FROM 
            comments c
        JOIN 
            users u ON c.user_id = u.id
        LEFT JOIN
            comment_votes cv ON c.id = cv.comment_id
        WHERE 
            c.link_id = $1
        GROUP BY
            c.id, u.username
		ORDER BY 
			vote_score DESC, c.created_at DESC
        LIMIT $2
        OFFSET $3
    `
	rows, err := q.db.Query(ctx, query, arg.LinkID, arg.Limit, arg.Offset, arg.UserID)
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
			&commentRow.ReplyCount,
			&commentRow.VoteScore,
			&commentRow.UserVoted,
			&commentRow.CreatedAt,
		); err != nil {
			return nil, err
		}
		commentRows = append(commentRows, commentRow)
	}
	return commentRows, nil
}

type CommentVoteParams struct {
	UserID    uuid.UUID
	CommentID uuid.UUID
	Vote      int32
}

func (q *Queries) CommentVote(ctx context.Context, arg CommentVoteParams) error {
	query := `
        WITH existing_vote AS (
            SELECT vote FROM comment_votes WHERE user_id = $1 AND comment_id = $2 FOR UPDATE
        ), deleted AS (
            DELETE FROM comment_votes WHERE user_id = $1 AND comment_id = $2 AND vote = $3
        ), updated AS (
            UPDATE comment_votes SET vote = $3 WHERE user_id = $1 AND comment_id = $2 AND vote != $3
        )
        INSERT INTO comment_votes (user_id, comment_id, vote)
        SELECT $1, $2, $3 WHERE NOT EXISTS (SELECT 1 FROM existing_vote);
    `
	_, err := q.db.Exec(ctx, query, arg.UserID, arg.CommentID, arg.Vote)
	return err
}
