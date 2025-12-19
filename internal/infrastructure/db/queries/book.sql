-- name: SaveBook :exec
INSERT INTO m_books (book_code, book_name)
VALUES ($1, $2)
ON CONFLICT (book_code) DO UPDATE 
SET book_name = EXCLUDED.book_name;

-- name: ListBooks :many
SELECT book_code, book_name 
FROM m_books 
ORDER BY book_code;