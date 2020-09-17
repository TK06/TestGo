package main

import (
	"TestGo/api"
	"TestGo/client"
	"TestGo/db"
)

func main() {

	// APIクライアントのテスト
	// 単純なGET
	client.TestAPIClientGet()
	// パラメーター付きGET
	client.TestAPIClientGetParam("1", "testname") // URLのパスパラメーター（仮）適当に設定
	// POST（パラメータがjson形式）
	client.TestAPIClientPost1()
	// POST（パラメータをURLエンコード）
	client.TestAPIClientPost2("https://httpbin.org/post", "99") // URL,id設定

	// DB接続（MySQL）
	db.InitDb()
	// DB切断（実行完了後DB接続を閉じる）
	defer db.CloseDb()

	// api
	api.API()
}
