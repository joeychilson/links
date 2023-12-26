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
    COUNT(DISTINCT l.id) AS vote_count
FROM 
    links l
JOIN 
    users u ON l.user_id = u.id
LEFT JOIN 
    comments c ON l.id = c.link_id
LEFT JOIN 
    votes v ON l.id = v.link_id
GROUP BY 
    l.id, u.username
ORDER BY 
    l.created_at DESC
LIMIT 
    $1
OFFSET 
    $2;

-- name: CreateVote :exec
INSERT INTO votes (user_id, link_id) VALUES ($1, $2);

-- name: CountVotes :one
SELECT COUNT(*) FROM votes WHERE link_id = $1;

-- name: UserVoted :one
SELECT EXISTS(SELECT 1 FROM votes WHERE user_id = $1 AND link_id = $2);