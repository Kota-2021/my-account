package excel

import (
	"my-account/internal/domain"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

// LoadJournalExcel は仕訳帳データを読み込みます
func LoadJournalExcel(filePath string) ([]domain.Journal, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	var journal []domain.Journal
	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue
		}

		date, _ := time.Parse("2006-01-02", row[2]) // C列: 日付
		// row[3]が¥と,で区切られている場合は、¥と,を取り除いて数値に変換する
		tmpWithdrawal := strings.ReplaceAll(strings.ReplaceAll(row[3], "¥", ""), ",", "")
		withdrawal, err := strconv.ParseFloat(tmpWithdrawal, 64) // D列: 支払
		if err != nil {
			return nil, err
		}
		// row[5]が￥と,で区切られている場合は、￥を取り除いて数値に変換する
		tmpDeposit := strings.ReplaceAll(strings.ReplaceAll(row[5], "¥", ""), ",", "")
		deposit, err := strconv.ParseFloat(tmpDeposit, 64) // F列: 入金
		if err != nil {
			return nil, err
		}
		subjectCode, err := strconv.Atoi(row[11]) // Q列: 科目コード
		if err != nil {
			return nil, err
		}
		item := row[13]     // R列: 摘要
		customer := row[14] // S列: 取引先
		evidence := ""
		if len(row) > 15 {
			evidence = row[15] // T列: 証票番号
		}
		memo := ""
		if len(row) > 16 {
			memo = row[16] // U列: 摘要
		}
		if err != nil {
			return nil, err
		}
		categoryID, err := strconv.Atoi(row[10]) // W列: カテゴリーID
		if err != nil {
			return nil, err
		}

		journal = append(journal, domain.Journal{
			Date:        date,
			Withdrawal:  decimal.NewFromFloat(withdrawal),
			Deposit:     decimal.NewFromFloat(deposit),
			SubjectCode: int16(subjectCode),
			Item:        item,
			Customer:    customer,
			Evidence:    evidence,
			Memo:        memo,
			// BookCode:    int16(bookCode),
			CategoryID: int16(categoryID),
			// FiscalYear:  int16(fiscalYear),
		})
	}
	return journal, nil
}
