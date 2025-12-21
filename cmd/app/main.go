package main

import (
	"context"
	"fmt"
	"log"
	"my-account/internal/infrastructure/db"
	"my-account/internal/infrastructure/excel"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {

	// 初期設定値
	const currentYear int16 = 2025
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current path")
	}
	envFilePath := currentPath + "/.env"
	const inputExcelBasePath string = "internal/infrastructure/excel/inputdata/"

	// .envファイルから環境変数を読み込む
	err = godotenv.Load(envFilePath)
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

	// DB接続情報を取得
	ctx := context.Background()
	// DB接続
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("接続失敗: %v", err)
	}
	defer conn.Close(ctx)

	// --- カテゴリーマスタ (m_categories) の処理 ---

	fmt.Println(">>> カテゴリーマスタの処理を開始")
	// Excel読込
	categories, err := excel.LoadCategoriesExcel(inputExcelBasePath + "categories.xlsx")
	if err != nil {
		log.Printf("カテゴリーExcel読込エラー: %v", err)
	} else {
		// トランザクション開始
		tx, err := conn.Begin(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback(ctx)

		// DB保存 (sqlcを利用)
		if err := db.SaveCategories(ctx, tx, categories); err != nil {
			log.Fatalf("カテゴリー保存エラー: %v", err)
		}
		tx.Commit(ctx)
		fmt.Println("カテゴリーマスタの更新が完了しました。")
	}

	// --- 勘定科目マスタ (m_subjects) の処理 ---

	fmt.Println("\n>>> 勘定科目マスタの処理を開始")
	// Excel読込
	subjects, err := excel.LoadSubjectsExcel(inputExcelBasePath + "subjects.xlsx")
	if err != nil {
		log.Printf("勘定科目Excel読込エラー: %v", err)
	} else {
		tx, err := conn.Begin(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback(ctx)

		// DB保存 (sqlcを利用)
		if err := db.SaveSubjects(ctx, tx, subjects); err != nil {
			log.Fatalf("勘定科目保存エラー: %v", err)
		}
		tx.Commit(ctx)
		fmt.Println("勘定科目マスタの更新が完了しました。")
	}

	// --- 帳票マスタ (m_books) の処理 ---

	fmt.Println("\n>>> 帳票マスタの処理を開始")
	// Excel読込
	books, err := excel.LoadBooksExcel(inputExcelBasePath + "books.xlsx")
	if err != nil {
		log.Printf("帳票マスタExcel読込エラー: %v", err)
	} else {
		tx, err := conn.Begin(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback(ctx)

		// DB保存 (sqlcを利用)
		if err := db.SaveBooks(ctx, tx, books); err != nil {
			log.Fatalf("帳票マスタ保存エラー: %v", err)
		}
		tx.Commit(ctx)
		fmt.Println("帳票マスタの更新が完了しました。")
	}

	// --- 予算データ (t_buget_financial_data) の処理 ---

	fmt.Println("\n>>> 予算データの処理を開始")
	// Excel読込
	bugets, err := excel.LoadBugetsExcel(inputExcelBasePath + "bugets.xlsx")
	if err != nil {
		log.Printf("予算データExcel読込エラー: %v", err)
	} else {
		tx, err := conn.Begin(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback(ctx)

		// DB保存 (sqlcを利用)
		if err := db.SaveBugets(ctx, tx, bugets); err != nil {
			log.Fatalf("予算データ保存エラー: %v", err)
		}
		tx.Commit(ctx)
		fmt.Println("予算データの更新が完了しました。")
	}

	// --- 出納帳データ (t_cashbook) の処理 ---

	// 帳票マスタ一覧を取得する。
	bookList, _ := db.FetchAllBooks(ctx, conn)
	fmt.Printf("[帳票マスタ] より %d 件を取得。\n", len(bookList))
	for _, b := range bookList {
		fmt.Printf("  ID:%d - %s\n", b.Code, b.Name)
	}

	fmt.Println("\n>>> 出納帳データの処理を開始")
	for _, b := range bookList {
		excelPath := inputExcelBasePath + "cashbook_" + b.Name + ".xlsx"
		// Excel読込
		fmt.Printf("  Excelパス: %s\n", excelPath)
		cashbooks, err := excel.LoadCashbooksExcel(excelPath)
		if err != nil {
			log.Printf("出納帳データExcel読込エラー: %v", err)
		} else {
			tx, err := conn.Begin(ctx)
			if err != nil {
				log.Fatal(err)
			}
			defer tx.Rollback(ctx)

			// DB保存 (sqlcを利用)
			if err := db.SaveCashbooks(ctx, tx, cashbooks, currentYear, b.Code); err != nil {
				log.Fatalf("出納帳データ保存エラー: %v", err)
			}
			tx.Commit(ctx)
			fmt.Println("出納帳データの更新が完了しました。")
		}
	}

	// --- 登録結果の表示 ---

	fmt.Println("\n--- 登録済みデータの確認 ---")

	// カテゴリー一覧表示
	savedCats, _ := db.FetchAllCategories(ctx, conn)
	fmt.Printf("[カテゴリー] %d件登録済み\n", len(savedCats))
	for _, c := range savedCats {
		fmt.Printf("  ID:%d - %s\n", c.ID, c.Name)
	}

	// 勘定科目一覧表示
	savedSubjects, _ := db.FetchAllSubjects(ctx, conn)
	fmt.Printf("[勘定科目] %d件登録済み\n", len(savedSubjects))
	for _, s := range savedSubjects {
		fmt.Printf("  ID:%d - %s\n", s.Code, s.Name)
	}

	// 帳票マスタ一覧表示
	savedBooks, _ := db.FetchAllBooks(ctx, conn)
	fmt.Printf("[帳票マスタ] %d件登録済み\n", len(savedBooks))
	for _, b := range savedBooks {
		fmt.Printf("  ID:%d - %s\n", b.Code, b.Name)
	}

	// 予算データ一覧表示
	savedBugets, _ := db.FetchAllBugets(ctx, conn)
	fmt.Printf("[予算データ] %d件登録済み\n", len(savedBugets))
	for _, b := range savedBugets {
		fmt.Printf("  ID:%d - 科目コード:%d - カテゴリーID:%d - 予算:%s - 実績:%s - 差異:%s - 年度:%d\n", b.ID, b.SubjectCode, b.CategoryID, b.Budget, b.Result, b.Difference, b.FiscalYear)
	}

	// 出納帳データ一覧表示
	savedCashbooks, _ := db.FetchAllCashbooks(ctx, conn)
	fmt.Printf("[出納帳データ] %d件登録済み\n", len(savedCashbooks))
	for _, c := range savedCashbooks {
		fmt.Printf("  ID:%d - 日付:%s - 摘要:%s - 支払:%s - 入金:%s - 残高:%s - 備考:%s - 帳票コード:%d - 年度:%d\n", c.ID, c.Date, c.Item, c.Withdrawal, c.Deposit, c.Balance, c.Remarks, c.BookCode, c.BookYear)
	}

	fmt.Println("✅処理が完了しました。")
}
