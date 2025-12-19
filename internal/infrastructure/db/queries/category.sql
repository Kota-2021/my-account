-- name: SaveCategory :exec
INSERT INTO m_categories (category_name)
VALUES ($1)
ON CONFLICT DO NOTHING;

-- name: ListCategories :many
SELECT category_id, category_name
FROM m_categories
ORDER BY category_id;