package excel

import (
	"my-account/internal/domain"

	"github.com/xuri/excelize/v2"
)

// LoadCategoriesExcel はExcelからカテゴリー名を読み込みます
func LoadCategoriesExcel(filePath string) ([]domain.Category, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	var categories []domain.Category
	for i, row := range rows {
		if i == 0 || len(row) == 0 { // ヘッダーをスキップ
			continue
		}
		categories = append(categories, domain.Category{ID: int16(i), Name: row[1]})
	}
	return categories, nil
}
