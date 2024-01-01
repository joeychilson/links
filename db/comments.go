package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

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
	query := `
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
	row := q.db.QueryRow(ctx, query, arg.ID, arg.UserID)
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

type PopularCommentsParams struct {
	Slug   string
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) PopularComments(ctx context.Context, arg PopularCommentsParams) ([]CommentRow, error) {
	query := `
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
		ORDER BY score DESC, replies DESC, created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := q.db.Query(ctx, query,
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

type ControversialCommentsParams struct {
	Slug   string
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) ControversialComments(ctx context.Context, arg ControversialCommentsParams) ([]CommentRow, error) {
	query := `
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
		ORDER BY replies DESC, score DESC, created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := q.db.Query(ctx, query,
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

type LatestCommentsParams struct {
	Slug   string
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) LatestComments(ctx context.Context, arg LatestCommentsParams) ([]CommentRow, error) {
	query := `
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
		ORDER BY created_at DESC, score DESC, replies DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := q.db.Query(ctx, query,
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

type CreateCommentParams struct {
	UserID  uuid.UUID
	LinkID  uuid.UUID
	Content string
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) error {
	query := "INSERT INTO comments (user_id, link_id, content) VALUES ($1, $2, $3);"
	_, err := q.db.Exec(ctx, query,
		arg.UserID,
		arg.LinkID,
		arg.Content,
	)
	return err
}

type CreateReplyParams struct {
	UserID   uuid.UUID
	LinkID   uuid.UUID
	ParentID uuid.UUID
	Content  string
}

func (q *Queries) CreateReply(ctx context.Context, arg CreateReplyParams) error {
	query := "INSERT INTO comments (user_id, link_id, parent_id, content) VALUES ($1, $2, $3, $4);"
	_, err := q.db.Exec(ctx, query,
		arg.UserID,
		arg.LinkID,
		arg.ParentID,
		arg.Content,
	)
	return err
}

type CreateVoteParams struct {
	UserID    uuid.UUID
	CommentID uuid.UUID
	Vote      int32
}

func (q *Queries) CreateVote(ctx context.Context, arg CreateVoteParams) error {
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

type Comment struct {
	ID     uuid.UUID
	LinkID uuid.UUID
}

func (q *Queries) CommentList(ctx context.Context) ([]Comment, error) {
	query := "SELECT id, link_id FROM comments"
	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Comment
	for rows.Next() {
		var i Comment
		if err := rows.Scan(&i.ID, &i.LinkID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
