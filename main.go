package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQLドライバ

	// "github.com/gorilla/mux"
	"github.com/gin-gonic/gin"
)

func main() {

	// Golang：ローカル、またはDocker内から実行される。
	// PostgreSQL：Docker内で動作。

	// golangのローカル接続であればtrue。falseの場合はDocker内から実行される。
	local := false
	var connStr string
	var appPort string

	if local {
		fmt.Println("ローカルからの実行です")
		// ローカルの場合は.envファイルから環境変数を読み込む。
		// ※※※.envファイルは.gitignoreとして設定しているため、gitで管理しない。※※※
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else {
		fmt.Println("docker-compose.ymlからの実行です")
	}

	// ローカルの場合は.envファイルから環境変数を読み込む。
	// docker-compose.ymlの場合はdocker-compose.ymlから環境変数を読み込む。
	dbHost := ""
	if local {
		dbHost = os.Getenv("DB_LOCAL_HOST")
	} else {
		dbHost = os.Getenv("DB_HOST")
	}
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	appPort = os.Getenv("APP_PORT")

	// 接続文字列を作成
	if local {
		connStr = fmt.Sprintf(
			"host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbUser, dbPass, dbName,
		)
		fmt.Println("connStr: ", connStr)
	} else {
		connStr = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbUser, dbPass, dbName,
		)

	}

	// データベースへの接続を試行 (リトライ処理)
	var db *sql.DB
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		log.Printf("DB接続を試行中... (試行 %d/%d)", i+1, maxRetries)

		// データベースを開く
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("sql.Openエラー: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// データベースにPingを送信して接続を確認
		err = db.Ping()
		if err == nil {
			log.Println("✅ データベースへの接続に成功しました！")
			break
		}

		log.Printf("db.Pingエラー: %v", err)
		db.Close() // 接続失敗時は閉じる
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ データベースへの接続に失敗しました: %v", err)
	}

	defer db.Close()

	// 接続成功後のアプリケーション処理（今回は省略）
	log.Println("アプリケーションが正常に起動し、動作中です。")

	// 以下から待ち受けを開始する
	// ローカルの場合はginを使用して待ち受けを開始する。
	// Dockerの場合はアプリコンテナを起動状態にする。
	if local {
		fmt.Println("ローカルからの実行です")
		// 簡易的なハンドラ関数
		r := gin.Default()
		r.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello! Application is running and connected to PostgreSQL."})
		})
		r.Run(":" + appPort)

	} else {
		fmt.Println("docker-compose.ymlからの実行です")
		// アプリコンテナを起動状態に保つための簡易的な処理
		select {}
	}
}
