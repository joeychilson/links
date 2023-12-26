-- name: CreateArticle :one
INSERT INTO articles (user_id, title, link) VALUES ($1, $2, $3) RETURNING *;

-- name: ArticleFeed :many
SELECT 
    a.id AS article_id,
    a.title,
    a.link,
    a.created_at,
    u.username,
    COUNT(c.id) AS comment_count
FROM 
    articles a
JOIN 
    users u ON a.user_id = u.id
LEFT JOIN 
    comments c ON a.id = c.article_id
GROUP BY 
    a.id, u.username
ORDER BY 
    a.created_at DESC
LIMIT 
    $1
OFFSET 
    $2;