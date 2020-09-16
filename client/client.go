package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// TestAPIClientGet 　net/httpパッケージを使用
func TestAPIClientGet() {
	// A simple HTTP Request & Response Service.「https://httpbin.org」を使用
	// 　URLに対してパラメータを追加したリクエストを送信すると、
	// 　クライアントのIPアドレスやユーザーエージェント、送信されたパラメータなどの情報がJSONデータとして返してくれる

	// httpのクライアントを作成する
	client := new(http.Client)

	// リクエストを作成（Get）
	// 第一引数: http.MethodGet は "GET" を指定するのと同じ
	// 第二引数: リクエストするURL
	// 第三引数: リクエストボディで フォーム、ファイル などのデータを送信する場合に使う。何も送信しない場合は nil を指定する
	req, err := http.NewRequest("GET", "https://httpbin.org/get", nil)
	if err != nil {
		panic(err.Error())
	}

	// リクエストを実行
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	// 関数を抜ける際に必ずresponseをcloseするようにdeferでcloseを呼ぶ
	defer res.Body.Close()

	// レスポンスの読み込み（レスポンスボディをすべて読み出す）　レスポンスの取得
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	// ステータスコード
	statusCode := res.StatusCode
	// ヘッダーを取得
	res.Header.Get("Content-Type")
	// res.Bodyの大きさ len(body) と同じ
	contentLength := res.ContentLength
	// リクエストURL
	reqURL := res.Request.URL.String()

	// レスポンスの表示（body は []byte バイト配列）
	fmt.Println("--testAPIClient--Result-Start-------------------------------------------")
	fmt.Println("StatusCode：" + strconv.Itoa(statusCode))                // 数値を文字列に変換
	fmt.Println("ContentLength：" + strconv.FormatInt(contentLength, 10)) // 数値int64を文字列に変換　引数設定10は10進数
	fmt.Println("RequestURL：" + reqURL)
	fmt.Println("Body：")
	fmt.Printf("%s", body)
	fmt.Println("--testAPIClient--Result-End-------------------------------------------")
}
