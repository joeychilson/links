-- name: CreateComment :one
INSERT INTO comments (user_id, link_id, content) VALUES ($1, $2, $3) RETURNING id;

-- name: CreateReply :one
INSERT INTO comments (user_id, link_id, parent_id, content) VALUES ($1, $2, $3, $4) RETURNING id;

-- name: Comment :one
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
    c.id = $1;

-- name: CommentFeed :many
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