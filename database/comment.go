package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateCommentParams struct {
	UserID   uuid.UUID
	LinkID   uuid.UUID
	ParentID uuid.NullUUID
	Content  string
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) error {
	query := `
		INSERT INTO comments (user_id, link_id, parent_id, content)
		VALUES ($1, $2, $3, $4)
	`
	_, err := q.db.Exec(ctx, query, arg.UserID, arg.LinkID, arg.ParentID, arg.Content)
	return err
}

type CommentRow struct {
	ID        uuid.UUID
	LinkID    uuid.UUID
	ParentID  uuid.UUID
	Username  string
	Content   string
	Replies   int64
	Score     int64
	UserVote  int32
	CreatedAt time.Time
	UpdatedAt time.Time
	Children  []CommentRow
}

type CommentFeedParams struct {
	LinkID uuid.UUID
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) CommentFeed(ctx context.Context, arg CommentFeedParams) ([]CommentRow, error) {
	query := `
		WITH RECURSIVE comment_tree AS (
			SELECT 
				c.id,
				c.link_id,
				c.parent_id,
				c.content,
				c.created_at,
				c.updated_at,
				u.username,
				(SELECT COUNT(*) FROM comments WHERE parent_id = c.id) AS replies,
				(SELECT COALESCE(SUM(vote), 0) FROM comment_votes WHERE comment_id = c.id) AS score,
				COALESCE(cv.vote, 0) AS user_vote
			FROM 
				comments c
			JOIN 
				users u ON c.user_id = u.id
			LEFT JOIN 
				comment_votes cv ON c.id = cv.comment_id AND cv.user_id = $2
			WHERE 
				c.link_id = $1 AND c.parent_id IS NULL
		
			UNION ALL
		
			SELECT 
				c.id,
				c.link_id,
				c.parent_id,
				c.content,
				c.created_at,
				c.updated_at,
				u.username,
				(SELECT COUNT(*) FROM comments WHERE parent_id = c.id) AS replies,
				(SELECT COALESCE(SUM(vote), 0) FROM comment_votes WHERE comment_id = c.id) AS score,
				COALESCE(cv.vote, 0) AS user_vote
			FROM 
				comments c
			JOIN 
				comment_tree ct ON c.parent_id = ct.id
			JOIN 
				users u ON c.user_id = u.id
			LEFT JOIN 
				comment_votes cv ON c.id = cv.comment_id AND cv.user_id = $2
		)
		SELECT * FROM comment_tree
		ORDER BY score DESC, created_at
		LIMIT $3 OFFSET $4;
	`
	rows, err := q.db.Query(ctx, query, arg.LinkID, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commentRows []CommentRow
	for rows.Next() {
		var commentRow CommentRow
		if err := rows.Scan(
			&commentRow.ID,
			&commentRow.LinkID,
			&commentRow.ParentID,
			&commentRow.Content,
			&commentRow.CreatedAt,
			&commentRow.UpdatedAt,
			&commentRow.Username,
			&commentRow.Replies,
			&commentRow.Score,
			&commentRow.UserVote,
		); err != nil {
			return nil, err
		}
		commentRows = append(commentRows, commentRow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return buildCommentTree(commentRows, uuid.Nil), nil
}

func buildCommentTree(comments []CommentRow, parentID uuid.UUID) []CommentRow {
	var tree []CommentRow
	for _, comment := range comments {
		if comment.ParentID == parentID {
			children := buildCommentTree(comments, comment.ID)
			comment.Children = children
			tree = append(tree, comment)
		}
	}
	return tree
}

type CommentParams struct {
	UserID    uuid.UUID
	CommentID uuid.UUID
}

func (q *Queries) Comment(ctx context.Context, arg CommentParams) (CommentRow, error) {
	query := `
        SELECT 
            c.id,
            c.link_id,
            c.parent_id,
            c.content,
            c.created_at,
            c.updated_at,
            u.username,
            (SELECT COUNT(*) FROM comments WHERE parent_id = c.id) AS replies,
            (SELECT COALESCE(SUM(vote), 0) FROM comment_votes WHERE comment_id = c.id) AS score,
            COALESCE((SELECT cv.vote FROM comment_votes cv WHERE cv.comment_id = c.id AND cv.user_id = $1), 0) AS user_vote
        FROM 
            comments c
        JOIN 
            users u ON c.user_id = u.id
        WHERE 
            c.id = $2
        GROUP BY
            c.id, u.username, c.link_id
    `
	var commentRow CommentRow
	if err := q.db.QueryRow(ctx, query, arg.UserID, arg.CommentID).Scan(
		&commentRow.ID,
		&commentRow.LinkID,
		&commentRow.ParentID,
		&commentRow.Content,
		&commentRow.CreatedAt,
		&commentRow.UpdatedAt,
		&commentRow.Username,
		&commentRow.Replies,
		&commentRow.Score,
		&commentRow.UserVote,
	); err != nil {
		return CommentRow{}, err
	}
	return commentRow, nil
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

type CommentScoreAndUserVoteParams struct {
	UserID    uuid.UUID
	CommentID uuid.UUID
}

type CommentScoreAndUserVoteRow struct {
	Score    int64
	UserVote int32
}

func (q *Queries) CommentScoreAndUserVote(ctx context.Context, args CommentScoreAndUserVoteParams) (CommentScoreAndUserVoteRow, error) {
	query := `
		SELECT 
			COALESCE(SUM(vote), 0) AS score,
			COALESCE((SELECT vote FROM comment_votes WHERE user_id = $1 AND comment_id = $2), 0) AS user_vote
		FROM 
			comment_votes
		WHERE 
			comment_id = $2
	`
	var row CommentScoreAndUserVoteRow
	if err := q.db.QueryRow(ctx, query, args.UserID, args.CommentID).Scan(&row.Score, &row.UserVote); err != nil {
		return CommentScoreAndUserVoteRow{}, err
	}
	return row, nil
}
