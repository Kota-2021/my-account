package excel

import (
	"my-account/internal/domain"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func LoadBooksExcel(filePath string) ([]domain.Book, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	var books []domain.Book
	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue
		}

		code, _ := strconv.Atoi(row[0]) // A列: コード
		books = append(books, domain.Book{
			Code: int16(code),
			Name: row[1], // B列: 帳票名
		})
	}
	return books, nil
}
