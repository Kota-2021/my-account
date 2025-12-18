package domain

import "github.com/shopspring/decimal"

// FinancialData 決算データ (t_financial_data)
type FinancialData struct {
	SubjectCode int             `json:"subject_code"`
	CategoryID  int             `json:"category_id"`
	Budget      decimal.Decimal `json:"budget"`     // 予算
	Result      decimal.Decimal `json:"result"`     // 実績 (仕訳+未収+未払の合算)
	Difference  decimal.Decimal `json:"difference"` // 差異
}
