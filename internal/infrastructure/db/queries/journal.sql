-- name: SaveJournal :exec
INSERT INTO t_journal
(journal_date, withdrawal, deposit, subject_code, item, customer, evidence, memo, book_code, category_id, fiscal_year)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: ListJournal :many
SELECT journal_id, journal_date, withdrawal, deposit, subject_code, item, customer, evidence, memo, book_code, category_id, fiscal_year
FROM t_journal
ORDER BY journal_id;
