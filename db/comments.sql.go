package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const comment = `-- name: Comment :one
SELECT 
    c.id,
    l.slug AS link_slug,
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
LEFT JOIN 
    links l ON c.link_id = l.id
WHERE 
    c.id = $1
`

type CommentParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

type CommentRow struct {
	ID        uuid.UUID
	LinkSlug  string
	ParentID  uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Replies   int64
	Score     int64
	UserVote  int16
	Children  []CommentRow
}

func (q *Queries) Comment(ctx context.Context, arg CommentParams) (CommentRow, error) {
	row := q.db.QueryRow(ctx, comment, arg.ID, arg.UserID)
	var i CommentRow
	err := row.Scan(
		&i.ID,
		&i.LinkSlug,
		&i.ParentID,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Replies,
		&i.Score,
		&i.UserVote,
	)
	return i, err
}

const commentFeed = `-- name: CommentFeed :many
WITH RECURSIVE comment_tree AS (
    SELECT 
        c.id,
        l.slug AS link_slug,
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
    LEFT JOIN 
        links l ON c.link_id = l.id
    WHERE 
        l.slug = $1 AND c.parent_id IS NULL

    UNION ALL

    SELECT 
        c.id,
        l.slug AS link_slug,
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
    LEFT JOIN 
        links l ON c.link_id = l.id
)
SELECT id, link_slug, parent_id, content, created_at, updated_at, username, replies, score, user_vote FROM comment_tree
ORDER BY score DESC, created_at DESC
LIMIT $3 OFFSET $4
`

type CommentFeedParams struct {
	Slug   string
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) CommentFeed(ctx context.Context, arg CommentFeedParams) ([]CommentRow, error) {
	rows, err := q.db.Query(ctx, commentFeed,
		arg.Slug,
		arg.UserID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CommentRow
	for rows.Next() {
		var i CommentRow
		if err := rows.Scan(
			&i.ID,
			&i.LinkSlug,
			&i.ParentID,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Username,
			&i.Replies,
			&i.Score,
			&i.UserVote,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return buildCommentTree(items, uuid.Nil), nil
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

const createComment = `-- name: CreateComment :one
INSERT INTO comments (user_id, link_id, content) VALUES ($1, $2, $3) RETURNING id
`

type CreateCommentParams struct {
	UserID  uuid.UUID
	LinkID  uuid.UUID
	Content string
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createComment, arg.UserID, arg.LinkID, arg.Content)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createReply = `-- name: CreateReply :one
INSERT INTO comments (user_id, link_id, parent_id, content) VALUES ($1, $2, $3, $4) RETURNING id
`

type CreateReplyParams struct {
	UserID   uuid.UUID
	LinkID   uuid.UUID
	ParentID uuid.UUID
	Content  string
}

func (q *Queries) CreateReply(ctx context.Context, arg CreateReplyParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createReply,
		arg.UserID,
		arg.LinkID,
		arg.ParentID,
		arg.Content,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
