// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: feeds.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const controversialFeed = `-- name: ControversialFeed :many
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
GROUP BY 
    l.id, u.username, ul.liked, c.comments, ll.likes
ORDER BY 
    comments DESC, likes DESC, l.created_at DESC
LIMIT
    $2
OFFSET
    $3
`

type ControversialFeedParams struct {
	Column1 uuid.UUID
	Limit   int32
	Offset  int32
}

type ControversialFeedRow struct {
	ID        uuid.UUID
	Title     string
	Url       string
	Slug      string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Username  string
	Comments  int64
	Likes     int64
	Liked     bool
}

func (q *Queries) ControversialFeed(ctx context.Context, arg ControversialFeedParams) ([]ControversialFeedRow, error) {
	rows, err := q.db.Query(ctx, controversialFeed, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ControversialFeedRow
	for rows.Next() {
		var i ControversialFeedRow
		if err := rows.Scan(
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

const latestFeed = `-- name: LatestFeed :many
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
GROUP BY 
    l.id, u.username, ul.liked, c.comments, ll.likes
ORDER BY 
    l.created_at DESC, likes DESC, comments DESC
LIMIT
    $2
OFFSET
    $3
`

type LatestFeedParams struct {
	Column1 uuid.UUID
	Limit   int32
	Offset  int32
}

type LatestFeedRow struct {
	ID        uuid.UUID
	Title     string
	Url       string
	Slug      string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Username  string
	Comments  int64
	Likes     int64
	Liked     bool
}

func (q *Queries) LatestFeed(ctx context.Context, arg LatestFeedParams) ([]LatestFeedRow, error) {
	rows, err := q.db.Query(ctx, latestFeed, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LatestFeedRow
	for rows.Next() {
		var i LatestFeedRow
		if err := rows.Scan(
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

const popularFeed = `-- name: PopularFeed :many
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
GROUP BY 
    l.id, u.username, ul.liked, c.comments, ll.likes
ORDER BY 
    likes DESC, comments DESC, l.created_at DESC
LIMIT
    $2
OFFSET
    $3
`

type PopularFeedParams struct {
	Column1 uuid.UUID
	Limit   int32
	Offset  int32
}

type PopularFeedRow struct {
	ID        uuid.UUID
	Title     string
	Url       string
	Slug      string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Username  string
	Comments  int64
	Likes     int64
	Liked     bool
}

func (q *Queries) PopularFeed(ctx context.Context, arg PopularFeedParams) ([]PopularFeedRow, error) {
	rows, err := q.db.Query(ctx, popularFeed, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PopularFeedRow
	for rows.Next() {
		var i PopularFeedRow
		if err := rows.Scan(
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
