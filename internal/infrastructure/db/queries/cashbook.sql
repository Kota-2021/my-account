-- name: SaveCashbook :exec
INSERT INTO t_cashbook
(cashbook_date, item, withdrawal, deposit, balance, remarks, book_code, book_year)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: ListCashbook :many
SELECT cashbook_id, cashbook_date, item, withdrawal, deposit, balance, remarks, book_code, book_year
FROM t_cashbook
ORDER BY cashbook_id;
