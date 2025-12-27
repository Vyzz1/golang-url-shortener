
-- name: InsertClick :one
INSERT INTO clicks (url_id, ip_address, clicked_at, user_agent, referer, device_type, country)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;