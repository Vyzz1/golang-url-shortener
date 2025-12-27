-- db/queries/urls.sql

-- name: CreateURL :one
INSERT INTO urls (short_code, original_url)
VALUES ($1, $2)
RETURNING *;

-- name: GetURLByShortCode :one
SELECT * FROM urls
WHERE short_code = $1 AND is_active = true
LIMIT 1;

-- name: IncrementClickCount :exec
UPDATE urls
SET click_count = click_count + 1
WHERE short_code = $1;

-- name: ListURLs :many
SELECT * FROM urls
WHERE is_active = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetURLStats :one
SELECT 
    short_code,
    original_url,
    click_count,
    created_at
FROM urls
WHERE short_code = $1 AND is_active = true
LIMIT 1;

-- name: CheckShortCodeExists :one
SELECT EXISTS(
    SELECT 1 FROM urls WHERE short_code = $1
) AS exists;

-- name: DeactivateURL :exec
UPDATE urls
SET is_active = false
WHERE short_code = $1;

-- name: CountURLs :one
SELECT COUNT(*) AS url_count
FROM urls
WHERE is_active = true;