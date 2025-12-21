package excel

import (
	"my-account/internal/domain"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

// LoadCashbooksExcel は出納帳データを読み込みます
func LoadCashbooksExcel(filePath string) ([]domain.Cashbook, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	var cashbooks []domain.Cashbook
	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue
		}

		date, _ := time.Parse("2006-01-02", row[0])       // A列: 日付
		item := row[1]                                    // B列: 摘要
		withdrawal, err := strconv.ParseFloat(row[2], 64) // B列: 支払
		if err != nil {
			return nil, err
		}
		deposit, err := strconv.ParseFloat(row[3], 64) // C列: 入金
		if err != nil {
			return nil, err
		}
		balance, err := strconv.ParseFloat(row[4], 64) // D列: 残高
		if err != nil {
			return nil, err
		}
		remarks := ""
		if len(row) > 5 {
			remarks = row[5] // E列: 備考
		}

		cashbooks = append(cashbooks, domain.Cashbook{
			Date:       date,
			Item:       item,
			Withdrawal: decimal.NewFromFloat(withdrawal),
			Deposit:    decimal.NewFromFloat(deposit),
			Balance:    decimal.NewFromFloat(balance),
			Remarks:    remarks,
		})
	}
	return cashbooks, nil
}
