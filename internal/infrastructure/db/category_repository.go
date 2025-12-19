package db

import (
	"context"
	"my-account/internal/domain"

	"github.com/jackc/pgx/v5"
)

// SaveCategories はカテゴリーを一括登録します
func SaveCategories(ctx context.Context, tx pgx.Tx, categories []domain.Category) error {
	for _, cat := range categories {
		_, err := tx.Exec(ctx,
			"INSERT INTO m_categories (category_name) VALUES ($1) ON CONFLICT DO NOTHING",
			cat.Name) //
		if err != nil {
			return err
		}
	}
	return nil
}

// FetchAllCategories は全カテゴリーを取得します
func FetchAllCategories(ctx context.Context, conn *pgx.Conn) ([]domain.Category, error) {
	rows, err := conn.Query(ctx, "SELECT category_id, category_name FROM m_categories ORDER BY category_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	return results, nil
}
