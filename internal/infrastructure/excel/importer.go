package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func main() {
	// 1. ファイルを開く
	f, err := excelize.OpenFile("7-family-tyoubo1-251202.xlsx") // 読み込むファイル名
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 2. データの取得 (シート名は要件に合わせて変更してください)
	// 「仕訳帳」などのシートから全行を取得
	rows, err := f.GetRows("7-family-tyoubo1-251202")
	if err != nil {
		log.Fatal(err)
	}

	// 3. 1行目と2行目を表示
	fmt.Println("--- 読込データの確認 ---")
	for i, row := range rows {
		if i >= 2 { // 0-indexed なので 0, 1 行目のみ処理
			break
		}

		fmt.Printf("%d行目: ", i+1)
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}
