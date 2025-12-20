package db

import (
	"context"
	"my-account/internal/domain"
	"my-account/internal/infrastructure/db/dbgen"

	"github.com/jackc/pgx/v5"
)

// SaveCategories (sqlc + pgx/v5 対応)
func SaveCategories(ctx context.Context, tx pgx.Tx, categories []domain.Category) error {
	q := dbgen.New(tx) // sql_package: "pgx/v5" により tx が受け入れ可能になります
	for _, cat := range categories {
		// CategoryName は varchar(50) への修正を反映させておいてください
		if err := q.SaveCategory(ctx, dbgen.SaveCategoryParams{
			CategoryID:   cat.ID,
			CategoryName: cat.Name,
		}); err != nil {
			return err
		}
	}
	return nil
}

// FetchAllCategories (sqlc + pgx/v5 対応)
func FetchAllCategories(ctx context.Context, conn *pgx.Conn) ([]domain.Category, error) {
	q := dbgen.New(conn)
	dbCats, err := q.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	var results []domain.Category
	for _, c := range dbCats {
		results = append(results, domain.Category{
			ID:   int16(c.CategoryID),
			Name: c.CategoryName,
		})
	}
	return results, nil
}
