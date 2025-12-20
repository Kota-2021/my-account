package excel

import (
	"my-account/internal/domain"
	"strconv"

	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

// LoadBugetsExcel は予算データを読み込みます
func LoadBugetsExcel(filePath string) ([]domain.BugetFinancialData, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	var bugets []domain.BugetFinancialData
	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue
		}

		subjectCode, _ := strconv.Atoi(row[0])        // A列: 科目コード
		budget, err := strconv.ParseFloat(row[1], 64) // B列: 予算
		if err != nil {
			return nil, err
		}
		result, err := strconv.ParseFloat(row[2], 64) // C列: 実績
		if err != nil {
			return nil, err
		}
		difference, err := strconv.ParseFloat(row[3], 64) // D列: 差異
		if err != nil {
			return nil, err
		}
		categoryID, _ := strconv.Atoi(row[4]) // D列: カテゴリーID
		fiscalYear, _ := strconv.Atoi(row[5]) // E列: 年度

		bugets = append(bugets, domain.BugetFinancialData{
			SubjectCode: int16(subjectCode),
			Budget:      decimal.NewFromFloat(budget),
			Result:      decimal.NewFromFloat(result),
			Difference:  decimal.NewFromFloat(difference),
			CategoryID:  int16(categoryID),
			FiscalYear:  int16(fiscalYear),
		})
	}
	return bugets, nil
}
