package domain

import "github.com/shopspring/decimal"

// BugetFinancialData 予算・決算データ (t_buget_financial_data)
type BugetFinancialData struct {
	ID          int             `json:"buget_financial_data_id"`
	SubjectCode int16           `json:"subject_code"`
	CategoryID  int16           `json:"category_id"`
	Budget      decimal.Decimal `json:"budget"`     // 予算
	Result      decimal.Decimal `json:"result"`     // 実績 (仕訳+未収+未払の合算)
	Difference  decimal.Decimal `json:"difference"` // 差異
	FiscalYear  int16           `json:"fiscal_year"`
}
