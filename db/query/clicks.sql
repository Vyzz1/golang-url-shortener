
-- name: InsertClick :one
INSERT INTO clicks (url_id, ip_address, clicked_at, user_agent, referer, device_type, country)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetClicksByURLID :many

SELECT * FROM clicks
WHERE url_id = $1
ORDER BY clicked_at DESC
LIMIT $2 OFFSET $3;


-- name: CountClicksByURLID :one
SELECT COUNT(*) AS click_count
FROM clicks
WHERE url_id = $1;