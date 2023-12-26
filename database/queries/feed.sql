-- name: CreateArticle :exec
INSERT INTO articles (user_id, title, link) VALUES ($1, $2, $3);

-- name: ArticleFeed :many
SELECT 
    a.id AS article_id,
    a.title,
    a.link,
    a.created_at,
    u.username,
    COUNT(DISTINCT c.id) AS comment_count,
    COUNT(DISTINCT l.id) AS like_count
FROM 
    articles a
JOIN 
    users u ON a.user_id = u.id
LEFT JOIN 
    comments c ON a.id = c.article_id
LEFT JOIN 
    likes l ON a.id = l.article_id
GROUP BY 
    a.id, u.username
ORDER BY 
    a.created_at DESC
LIMIT 
    $1
OFFSET 
    $2;

-- name: CreateLike :exec
INSERT INTO likes (user_id, article_id) VALUES ($1, $2);

-- name: CountLikes :one
SELECT COUNT(*) FROM likes WHERE article_id = $1;