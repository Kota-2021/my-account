package db

import (
	"context"
	"log"
	"my-account/internal/domain"
	"my-account/internal/infrastructure/db/dbgen"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

// SaveBugets は予算データを一括登録・更新します
func SaveBugets(ctx context.Context, tx pgx.Tx, bugets []domain.BugetFinancialData) error {
	q := dbgen.New(tx)
	for _, b := range bugets {
		// shopspring/decimal を pgtype.Numeric に変換
		var budget, result, diff pgtype.Numeric
		budget.Scan(b.Budget.String()) // b.Budget が decimal.Decimal 型と仮定
		result.Scan(b.Result.String())
		diff.Scan(b.Difference.String())

		sbp := dbgen.SaveBugetParams{
			SubjectCode:     b.SubjectCode,
			Budget:          budget,
			Result:          result,
			Difference:      diff,
			CategoryID:      b.CategoryID,
			BugetFiscalYear: b.FiscalYear,
		}

		err := q.SaveBuget(ctx, sbp)
		if err != nil {
			log.Printf("予算データ保存エラー: %v", err)
			return err
		}
	}
	return nil
}

func FetchAllBugets(ctx context.Context, conn *pgx.Conn) ([]domain.BugetFinancialData, error) {
	q := dbgen.New(conn)
	dbBugets, err := q.ListBuget(ctx)
	if err != nil {
		return nil, err
	}

	var results []domain.BugetFinancialData
	for _, b := range dbBugets {
		results = append(results, domain.BugetFinancialData{
			ID:          int(b.BugetFinancialDataID),
			SubjectCode: b.SubjectCode,
			Budget:      decimal.NewFromBigInt(b.Budget.Int, b.Budget.Exp),
			Result:      decimal.NewFromBigInt(b.Result.Int, b.Result.Exp),
			Difference:  decimal.NewFromBigInt(b.Difference.Int, b.Difference.Exp),
			CategoryID:  b.CategoryID,
			FiscalYear:  b.BugetFiscalYear,
		})
	}
	return results, nil
}
