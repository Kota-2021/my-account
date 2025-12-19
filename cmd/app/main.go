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

	// 1. Excel読込
	fmt.Println("Excelを読み込み中...")
	cats, err := excel.LoadCategoriesExcel("categories.xlsx")
	if err != nil {
		log.Fatalf("読込エラー: %v", err)
	}

	// 2. DB書き込み (トランザクション)
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback(ctx)

	if err := db.SaveCategories(ctx, tx, cats); err != nil {
		log.Fatalf("保存エラー: %v", err)
	}
	tx.Commit(ctx)
	fmt.Println("DBへの書き込みが完了しました。")

	// 3. レコードの読みだし確認
	fmt.Println("\n--- 現在のカテゴリーマスタ ---")
	savedCats, err := db.FetchAllCategories(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range savedCats {
		fmt.Printf("[%d] %s\n", c.ID, c.Name)
	}
}
