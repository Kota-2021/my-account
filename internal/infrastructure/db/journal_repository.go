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

// SaveCashbooks は出納帳データを一括登録・更新します
func SaveJournal(ctx context.Context, tx pgx.Tx, journal []domain.Journal, currentYear int16, bookCode int16) error {
	q := dbgen.New(tx)
	for _, j := range journal {
		// shopspring/decimal を pgtype.Numeric に変換
		var withdrawal, deposit pgtype.Numeric
		withdrawal.Scan(j.Withdrawal.String())
		deposit.Scan(j.Deposit.String())

		sbp := dbgen.SaveJournalParams{
			JournalDate: pgtype.Date{Time: j.Date, Valid: true},
			Withdrawal:  withdrawal,
			Deposit:     deposit,
			SubjectCode: j.SubjectCode,
			Item:        pgtype.Text{String: j.Item, Valid: true},
			Customer:    j.Customer,
			Evidence:    pgtype.Text{String: j.Evidence, Valid: true},
			Memo:        pgtype.Text{String: j.Memo, Valid: true},
			BookCode:    bookCode,
			CategoryID:  j.CategoryID,
			FiscalYear:  currentYear,
		}

		err := q.SaveJournal(ctx, sbp)
		if err != nil {
			log.Printf("仕訳帳データ保存エラー: %v", err)
			return err
		}
	}
	return nil
}

func FetchAllJournal(ctx context.Context, conn *pgx.Conn) ([]domain.Journal, error) {
	q := dbgen.New(conn)
	dbJournal, err := q.ListJournal(ctx)
	if err != nil {
		return nil, err
	}

	var results []domain.Journal
	for _, j := range dbJournal {
		results = append(results, domain.Journal{
			ID:          int(j.JournalID),
			Date:        j.JournalDate.Time,
			Withdrawal:  decimal.NewFromBigInt(j.Withdrawal.Int, j.Withdrawal.Exp),
			Deposit:     decimal.NewFromBigInt(j.Deposit.Int, j.Deposit.Exp),
			SubjectCode: int16(j.SubjectCode),
			Item:        j.Item.String,
			Customer:    j.Customer,
			Evidence:    j.Evidence.String,
			Memo:        j.Memo.String,
			BookCode:    int16(j.BookCode),
			CategoryID:  int16(j.CategoryID),
			FiscalYear:  int16(j.FiscalYear),
		})
	}
	return results, nil
}
