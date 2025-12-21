package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// Journal 仕訳帳データ (t_journal)
type Journal struct {
	ID          int             `json:"journal_id"`
	Date        time.Time       `json:"journal_date"` // 実際の実装では time.Time
	Withdrawal  decimal.Decimal `json:"withdrawal"`   // 支払
	Deposit     decimal.Decimal `json:"deposit"`      // 入金
	SubjectCode int16           `json:"subject_code"`
	Item        string          `json:"item"`        // 摘要
	Customer    string          `json:"customer"`    // 取引先
	Evidence    string          `json:"evidence"`    // 証票番号
	Memo        string          `json:"memo"`        // 摘要
	BookCode    int16           `json:"book_code"`   // 帳票コード
	CategoryID  int16           `json:"category_id"` // 6つのカテゴリー紐付け
	FiscalYear  int16           `json:"fiscal_year"`
}

// Receivable 未収金データ (t_receivable)
type Receivable struct {
	ID          int             `json:"journal_id"`
	Date        string          `json:"journal_date"` // 実際の実装では time.Time
	Withdrawal  decimal.Decimal `json:"withdrawal"`   // 支払
	Deposit     decimal.Decimal `json:"deposit"`      // 入金
	SubjectCode int16           `json:"subject_code"`
	Item        string          `json:"item"`        // 摘要
	Customer    string          `json:"customer"`    // 取引先
	CategoryID  int16           `json:"category_id"` // 6つのカテゴリー紐付け
	FiscalYear  int16           `json:"fiscal_year"`
}

// Payable 未払金データ (t_payable)
type Payable struct {
	ID          int             `json:"journal_id"`
	Date        string          `json:"journal_date"` // 実際の実装では time.Time
	Withdrawal  decimal.Decimal `json:"withdrawal"`   // 支払
	Deposit     decimal.Decimal `json:"deposit"`      // 入金
	SubjectCode int16           `json:"subject_code"`
	Item        string          `json:"item"`        // 摘要
	Customer    string          `json:"customer"`    // 取引先
	CategoryID  int16           `json:"category_id"` // 6つのカテゴリー紐付け
	FiscalYear  int16           `json:"fiscal_year"`
}

// Cashbook 出納帳データ (t_cashbook)
type Cashbook struct {
	ID         int             `json:"cashbook_id"`
	Date       time.Time       `json:"cashbook_date"`
	Item       string          `json:"item"`
	Withdrawal decimal.Decimal `json:"withdrawal"`
	Deposit    decimal.Decimal `json:"deposit"`
	Balance    decimal.Decimal `json:"balance"`
	Remarks    string          `json:"remarks"`
	BookCode   int16           `json:"book_code"`
	BookYear   int16           `json:"book_year"`
}
