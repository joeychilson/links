-- name: CreateLink :exec
INSERT INTO links (user_id, title, url) VALUES ($1, $2, $3);

-- name: LinkFeed :many
SELECT 
    l.id AS id,
    l.title,
    l.url,
    l.created_at,
    u.username,
    COUNT(DISTINCT c.id) AS comment_count,
    COUNT(DISTINCT lk.id) AS like_count,
    SUM(CASE WHEN lk.user_id = $1 THEN 1 ELSE 0 END) AS user_voted
FROM 
    links l
JOIN 
    users u ON l.user_id = u.id
LEFT JOIN 
    comments c ON l.id = c.link_id
LEFT JOIN 
    likes lk ON l.id = lk.link_id
GROUP BY 
    l.id, u.username
ORDER BY 
    l.created_at DESC
LIMIT 
    $2
OFFSET 
    $3;

-- name: CreateLike :exec
INSERT INTO likes (user_id, link_id) VALUES ($1, $2);

-- name: DeleteVote :exec
DELETE FROM likes WHERE user_id = $1 AND link_id = $2;

-- name: CountLikes :one
SELECT COUNT(*) FROM likes WHERE link_id = $1;

-- name: UserLiked :one
SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND link_id = $2);