package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateLikeParams struct {
	UserID uuid.UUID
	LinkID uuid.UUID
}

func (q *Queries) CreateLike(ctx context.Context, arg CreateLikeParams) error {
	query := "INSERT INTO link_likes (user_id, link_id) VALUES ($1, $2)"
	_, err := q.db.Exec(ctx, query, arg.UserID, arg.LinkID)
	return err
}

type CreateLinkParams struct {
	UserID uuid.UUID
	Title  string
	Url    string
	Slug   string
}

func (q *Queries) CreateLink(ctx context.Context, arg CreateLinkParams) (string, error) {
	query := "INSERT INTO links (user_id, title, url, slug) VALUES ($1, $2, $3, $4) RETURNING slug"
	row := q.db.QueryRow(ctx, query,
		arg.UserID,
		arg.Title,
		arg.Url,
		arg.Slug,
	)
	var slug string
	err := row.Scan(&slug)
	return slug, err
}

type DeleteLikeParams struct {
	UserID uuid.UUID
	LinkID uuid.UUID
}

func (q *Queries) DeleteLike(ctx context.Context, arg DeleteLikeParams) error {
	query := "DELETE FROM link_likes WHERE user_id = $1 AND link_id = $2"
	_, err := q.db.Exec(ctx, query, arg.UserID, arg.LinkID)
	return err
}

type LinkRow struct {
	ID        uuid.UUID
	Title     string
	Url       string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Comments  int64
	Likes     int64
	Liked     bool
}

type LinkBySlugParams struct {
	UserID uuid.UUID
	Slug   string
}

func (q *Queries) LinkBySlug(ctx context.Context, arg LinkBySlugParams) (LinkRow, error) {
	query := `
		SELECT 
			l.id AS id,
			l.title,
			l.url,
			l.slug,
			l.created_at,
			l.updated_at,
			u.username,
			COALESCE(c.comments, 0) AS comments,
			COALESCE(ll.likes, 0) AS likes,
			COALESCE(ul.liked, FALSE) AS liked
		FROM
			links l
		JOIN
			users u ON l.user_id = u.id
		LEFT JOIN
			(SELECT
				link_id,
				COUNT(*) AS comments
			FROM
				comments
			GROUP BY link_id
			) c ON l.id = c.link_id
		LEFT JOIN
			(SELECT
				link_id,
				COUNT(*) AS likes
			FROM 
				link_likes
			GROUP BY link_id
			) ll ON l.id = ll.link_id
		LEFT JOIN
			(SELECT
				link_id,
				TRUE AS liked
			FROM
				link_likes
			WHERE
				user_id = $1::uuid
			) ul ON l.id = ul.link_id
		WHERE
			l.slug = $2
	`
	row := q.db.QueryRow(ctx, query, arg.UserID, arg.Slug)
	var i LinkRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Url,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Comments,
		&i.Likes,
		&i.Liked,
	)
	return i, err
}

func (q *Queries) LinkIDBySlug(ctx context.Context, slug string) (uuid.UUID, error) {
	query := "SELECT id FROM links WHERE slug = $1"
	row := q.db.QueryRow(ctx, query, slug)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

type LinkLikesAndLikedParams struct {
	LinkID uuid.UUID
	UserID uuid.UUID
}

type LinkLikesAndLikedRow struct {
	Likes int64
	Liked bool
}

func (q *Queries) LinkLikesAndLiked(ctx context.Context, arg LinkLikesAndLikedParams) (LinkLikesAndLikedRow, error) {
	query := `
		SELECT 
			COUNT(l.id) AS Likes,
			EXISTS (
				SELECT 1
				FROM link_likes ul
				WHERE ul.link_id = $1::uuid AND ul.user_id = $2::uuid
			) AS Liked
		FROM 
			link_likes l
		WHERE 
			l.link_id = $1::uuid
	`
	row := q.db.QueryRow(ctx, query, arg.LinkID, arg.UserID)
	var i LinkLikesAndLikedRow
	err := row.Scan(&i.Likes, &i.Liked)
	return i, err
}

func (q *Queries) CountLinks(ctx context.Context) (int64, error) {
	query := "SELECT COUNT(*) FROM links"
	row := q.db.QueryRow(ctx, query)
	var count int64
	err := row.Scan(&count)
	return count, err
}

type Link struct {
	ID uuid.UUID
}

func (q *Queries) LinkList(ctx context.Context) ([]Link, error) {
	query := "SELECT id FROM links"
	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Link{}
	for rows.Next() {
		var i Link
		if err := rows.Scan(&i.ID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
