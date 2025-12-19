package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

// Domain Model: フォルダ構成 [cite: 1] の internal/domain/journal.go に相当
type Journal struct {
	Date        time.Time
	Withdrawal  decimal.Decimal
	Deposit     decimal.Decimal
	SubjectCode int16
	Item        string
	Customer    string
	Evidence    string
	Memo        string
	BookCode    int16
	CategoryID  int32
	FiscalYear  int16
}

func main() {
	// 1. データベース接続設定 (システム構成 の PostgreSQL Port 5432 へ接続)
	// 実際には internal/config [cite: 3] から読み込む形式が望ましい

	// path print
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current path")
	}
	fmt.Println("current path: ", currentPath)
	// .envファイルから環境変数を読み込む
	err = godotenv.Load(currentPath + "/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_LOCAL_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbTZ := os.Getenv("TZ")

	// DB接続文字列を作成
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?timezone=%s", dbUser, dbPass, dbHost, dbPort, dbName, dbTZ)

	// DB接続を確立
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	fmt.Println("DB接続成功")

	// 2. サンプルデータの準備

	// カテゴリーマスタ (m_categories)
	categories := []string{"クラブ本体", "グッズ", "ファミリー"}

	// 仕訳データ (t_journal)
	sampleJournal := Journal{
		Date:        time.Now(),
		Withdrawal:  decimal.NewFromInt(10500), // shopspring/decimalを使用
		Deposit:     decimal.NewFromInt(0),
		SubjectCode: 101, // 現金など
		Item:        "事務用品代",
		Customer:    "B文房具店",
		Evidence:    "領収書",
		Memo:        "4月分消耗品",
		BookCode:    1,
		CategoryID:  1,    // クラブ本体
		FiscalYear:  2024, // 予算・決算は年次単位 [cite: 2]
	}

	// 3. データの書き込み実行

	// トランザクションの開始（会計システムの整合性確保のため）
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback(ctx)

	// 1. カテゴリーの登録
	for _, name := range categories {
		_, err := tx.Exec(ctx,
			"INSERT INTO m_categories (category_name) VALUES ($1) ON CONFLICT DO NOTHING",
			name)
		if err != nil {
			log.Fatalf("Failed to insert category: %v", err)
		}
	}

	// 2. 勘定科目マスタ (m_subjects) の登録
	_, err = tx.Exec(ctx,
		"INSERT INTO m_subjects (subject_code, subject_name) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		101, "現金")
	if err != nil {
		log.Fatalf("Failed to insert subject: %v", err)
	}

	// 3. 帳票マスタ (m_books) の登録
	_, err = tx.Exec(ctx,
		"INSERT INTO m_books (book_code, book_name) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		1, "現預金出納帳")
	if err != nil {
		log.Fatalf("Failed to insert book: %v", err)
	}

	// 4. 仕訳データの登録
	query := `
		INSERT INTO t_journal (
			journal_date, withdrawal, deposit, subject_code, item, 
			customer, evidence, memo, book_code, category_id, fiscal_year
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(ctx, query,
		sampleJournal.Date,
		sampleJournal.Withdrawal,
		sampleJournal.Deposit,
		sampleJournal.SubjectCode,
		sampleJournal.Item,
		sampleJournal.Customer,
		sampleJournal.Evidence,
		sampleJournal.Memo,
		sampleJournal.BookCode,
		sampleJournal.CategoryID,
		sampleJournal.FiscalYear,
	)
	if err != nil {
		log.Fatalf("Failed to insert journal: %v", err)
	}

	// コミット
	if err := tx.Commit(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("サンプルデータの書き込みが正常に完了しました。")

	// *****************************
	// 保存したレコードを取得する
	// *****************************
	// --- 追記：保存したレコードを閲覧するサンプル ---

	fmt.Println("\n--- 登録済み仕訳データの確認 ---")

	// 1. SELECTクエリの実行
	// 最新の5件を取得する例
	rows, err := conn.Query(ctx, `
		SELECT 
			journal_date, withdrawal, deposit, item, customer, memo 
		FROM t_journal 
		ORDER BY journal_id DESC 
		LIMIT 5
	`)
	if err != nil {
		log.Fatalf("Failed to query journals: %v", err)
	}
	defer rows.Close()

	// 2. 結果のスキャンと表示
	for rows.Next() {
		var (
			date       time.Time
			withdrawal decimal.Decimal
			deposit    decimal.Decimal
			item       string
			customer   string
			memo       string
		)

		err := rows.Scan(&date, &withdrawal, &deposit, &item, &customer, &memo)
		if err != nil {
			log.Fatalf("Row scan failed: %v", err)
		}

		// 会計データとして整形して表示
		fmt.Printf("日付: %s | 項目: %-10s | 支払: %10s | 収入: %10s | 相手先: %s | 備考: %s\n",
			date.Format("2006-01-02"),
			item,
			withdrawal.StringFixed(0), // 小数点を表示せず整数の文字列として出力
			deposit.StringFixed(0),
			customer,
			memo,
		)
	}

	// rows.Next() ループ中でのエラーチェック
	if err := rows.Err(); err != nil {
		log.Fatalf("Error during row iteration: %v", err)
	}

}
