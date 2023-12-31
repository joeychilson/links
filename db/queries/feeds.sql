-- name: PopularFeed :many
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
    $3;

-- name: LatestFeed :many
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
    $3;

-- name: ControversialFeed :many
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
    $3;
