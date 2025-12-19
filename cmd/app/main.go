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
	ctx := context.Background()

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

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("接続失敗: %v", err)
	}
	defer conn.Close(ctx)

	// --- カテゴリーマスタ (m_categories) の処理 ---

	fmt.Println(">>> カテゴリーマスタの処理を開始")
	// Excel読込
	categories, err := excel.LoadCategoriesExcel("categories.xlsx")
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
	subjects, err := excel.LoadSubjectsExcel("subjects.xlsx")
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
}
