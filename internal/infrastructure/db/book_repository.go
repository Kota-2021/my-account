package db

import (
	"context"
	"my-account/internal/domain"
	"my-account/internal/infrastructure/db/dbgen"

	"github.com/jackc/pgx/v5"
)

// 参考コード
// type SaveBookParams struct {
// 	BookCode int16  `json:"book_code"`
// 	BookName string `json:"book_name"`
// }

// func (q *Queries) SaveBook(ctx context.Context, arg SaveBookParams) error {
// 	_, err := q.db.Exec(ctx, saveBook, arg.BookCode, arg.BookName)
// 	return err
// }

// SaveBooks は帳票マスタを一括登録・更新します
func SaveBooks(ctx context.Context, tx pgx.Tx, books []domain.Book) error {
	q := dbgen.New(tx)
	for _, b := range books {
		// book_codeはsmallintなのでint16にキャスト
		err := q.SaveBook(ctx, dbgen.SaveBookParams{
			BookCode: int16(b.Code),
			BookName: b.Name,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func FetchAllBooks(ctx context.Context, conn *pgx.Conn) ([]domain.Book, error) {
	q := dbgen.New(conn)
	dbBooks, err := q.ListBooks(ctx)
	if err != nil {
		return nil, err
	}

	var results []domain.Book
	for _, b := range dbBooks {
		results = append(results, domain.Book{
			Code: int16(b.BookCode),
			Name: b.BookName,
		})
	}
	return results, nil
}
