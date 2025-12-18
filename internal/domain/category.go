package domain

// Category カテゴリーマスタ (m_categories)
type Category struct {
	ID   int    `json:"category_id"`
	Name string `json:"category_name"` // クラブ本体, グッズ等
}

// Subject 勘定科目マスタ (m_subjects)
type Subject struct {
	Code int    `json:"subject_code"`
	Name string `json:"subject_name"`
}
