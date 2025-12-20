-- name: SaveCategory :exec
INSERT INTO m_categories (category_id, category_name)
VALUES ($1, $2)
ON CONFLICT (category_id) DO UPDATE 
SET category_name = EXCLUDED.category_name;

-- name: ListCategories :many
SELECT category_id, category_name
FROM m_categories
ORDER BY category_id;
