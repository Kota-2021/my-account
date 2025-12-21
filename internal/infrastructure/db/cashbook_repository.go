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
func SaveCashbooks(ctx context.Context, tx pgx.Tx, cashbooks []domain.Cashbook, currentYear int16, bookCode int16) error {
	q := dbgen.New(tx)
	for _, c := range cashbooks {
		// test print
		log.Printf("id: %d", c.ID)
		log.Printf("cashbook_date: %s", c.Date)
		log.Printf("item: %s", c.Item)
		log.Printf("withdrawal: %s", c.Withdrawal.String())
		log.Printf("deposit: %s", c.Deposit.String())
		log.Printf("balance: %s", c.Balance.String())
		log.Printf("remarks: %s", c.Remarks)
		log.Printf("book_code: %d", bookCode)
		log.Printf("book_year: %d", currentYear)

		// shopspring/decimal を pgtype.Numeric に変換
		var withdrawal, deposit, balance pgtype.Numeric
		withdrawal.Scan(c.Withdrawal.String())
		deposit.Scan(c.Deposit.String())
		balance.Scan(c.Balance.String())

		sbp := dbgen.SaveCashbookParams{
			CashbookDate: pgtype.Date{Time: c.Date, Valid: true},
			Item:         pgtype.Text{String: c.Item, Valid: true},
			Withdrawal:   withdrawal,
			Deposit:      deposit,
			Balance:      balance,
			Remarks:      pgtype.Text{String: c.Remarks, Valid: true},
			BookCode:     bookCode,
			BookYear:     currentYear,
		}
		// test print
		log.Printf("sbp: %+v", sbp)

		err := q.SaveCashbook(ctx, sbp)
		if err != nil {
			log.Printf("出納帳データ保存エラー: %v", err)
			return err
		}
	}
	return nil
}

func FetchAllCashbooks(ctx context.Context, conn *pgx.Conn) ([]domain.Cashbook, error) {
	q := dbgen.New(conn)
	dbCashbooks, err := q.ListCashbook(ctx)
	if err != nil {
		return nil, err
	}

	var results []domain.Cashbook
	for _, c := range dbCashbooks {
		results = append(results, domain.Cashbook{
			ID:         int(c.CashbookID),
			Date:       c.CashbookDate.Time,
			Withdrawal: decimal.NewFromBigInt(c.Withdrawal.Int, c.Withdrawal.Exp),
			Deposit:    decimal.NewFromBigInt(c.Deposit.Int, c.Deposit.Exp),
			Balance:    decimal.NewFromBigInt(c.Balance.Int, c.Balance.Exp),
			Remarks:    c.Remarks.String,
			BookCode:   c.BookCode,
			BookYear:   c.BookYear,
		})
	}
	return results, nil
}
