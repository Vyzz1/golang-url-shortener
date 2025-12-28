
-- name: CountAllClicks :one
SELECT COUNT(*) FROM clicks;

-- name: CountURLsToday :one
SELECT COUNT(*) FROM urls
WHERE DATE(created_at) = CURRENT_DATE;

-- name: CountClicksToday :one
SELECT COUNT(*) FROM clicks
WHERE DATE(clicked_at) = CURRENT_DATE;

-- name: GetTopURLs :many
SELECT 
    u.short_code,
    u.original_url,
    u.click_count
FROM urls u
ORDER BY u.click_count DESC
LIMIT $1;