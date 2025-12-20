-- name: SaveBuget :exec
INSERT INTO t_buget_financial_data
(subject_code, budget, result, difference, category_id, buget_fiscal_year)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (subject_code, buget_fiscal_year) DO UPDATE 
SET subject_code = EXCLUDED.subject_code,
    budget = EXCLUDED.budget,
    result = EXCLUDED.result,
    difference = EXCLUDED.difference,
    category_id = EXCLUDED.category_id,
    buget_fiscal_year = EXCLUDED.buget_fiscal_year;

-- name: ListBuget :many
SELECT buget_financial_data_id, subject_code, budget, result, difference, category_id, buget_fiscal_year
FROM t_buget_financial_data
ORDER BY buget_financial_data_id;
