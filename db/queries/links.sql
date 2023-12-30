-- name: CreateLink :one
INSERT INTO links (user_id, title, url, slug) VALUES ($1, $2, $3, $4) RETURNING slug;

-- name: LinkBySlug :one
SELECT * FROM links WHERE slug = $1;