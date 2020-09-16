package main

import (
	"TestGo/api"
	"TestGo/client"
	"TestGo/db"
	// フレームワーク Gin
	// GORM Go言語のORマッパー
)

func main() {

	// APIクライアントのテスト
	// 単純なGET
	client.TestAPIClientGet()
	// パラメーター付きGET
	//testAPIClientGetParam()
	// POST
	//testAPIClientPost()

	// DB接続（MySQL）
	db.InitDb()
	// DB切断（実行完了後DB接続を閉じる）
	defer db.CloseDb()

	// api
	api.API()
}
