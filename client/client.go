package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// TestAPIClientGet 　net/httpパッケージを使用
func TestAPIClientGet() {
	// A simple HTTP Request & Response Service.「https://httpbin.org」を使用
	// 　URLに対してパラメータを追加したリクエストを送信すると、
	// 　クライアントのIPアドレスやユーザーエージェント、送信されたパラメータなどの情報がJSONデータとして返してくれる

	// httpのクライアントを作成する
	client := new(http.Client)

	// リクエストを作成（Get）
	// 第一引数: http.MethodGet は "GET" を指定するのと同じ。メソッドを指定する（GET,POST,PUT,DELETE）
	// 第二引数: リクエストするURL
	// 第三引数: リクエストボディで フォーム、ファイル などのデータを送信する場合に使う。何も送信しない場合は nil を指定する
	req, err := http.NewRequest("GET", "https://httpbin.org/get", nil) // 「https://httpbin.org」サイトのHTTPMethodsのGETを使用
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
	header := res.Header.Get("Content-Type")
	// res.Bodyの大きさ len(body) と同じ
	contentLength := res.ContentLength
	// リクエストURL
	reqURL := res.Request.URL.String()

	// レスポンスの表示（body は []byte バイト配列）
	fmt.Println("--TestAPIClientGet--Result-Start-------------------------------------------")
	fmt.Println("StatusCode：" + strconv.Itoa(statusCode)) // 数値を文字列に変換
	fmt.Println("Header Content-Type：" + header)
	fmt.Println("ContentLength：" + strconv.FormatInt(contentLength, 10)) // 数値int64を文字列に変換　引数設定10は10進数
	fmt.Println("RequestURL：" + reqURL)
	fmt.Println("Body：")
	fmt.Printf("%s", body)
	fmt.Println("--TestAPIClientGet--Result-End-------------------------------------------")
}

// TestAPIClientGetParam 　GET パラメーター付き
func TestAPIClientGetParam(code string, name string) {
	// A simple HTTP Request & Response Service.「https://httpbin.org」を使用
	// 　URLに対してパラメータを追加したリクエストを送信すると、
	// 　クライアントのIPアドレスやユーザーエージェント、送信されたパラメータなどの情報がJSONデータとして返してくれる

	// httpのクライアントを作成する
	client := new(http.Client)

	// リクエストを作成（Get）
	// 第一引数: http.MethodGet は "GET" を指定するのと同じ。メソッドを指定する（GET,POST,PUT,DELETE）
	// 第二引数: リクエストするURL
	// 第三引数: リクエストボディで フォーム、ファイル などのデータを送信する場合に使う。何も送信しない場合は nil を指定する
	req, err := http.NewRequest("GET", "https://httpbin.org/get?code="+code+"&name="+name, nil) // 「https://httpbin.org」サイトのHTTPMethodsのGETを使用
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
	header := res.Header.Get("Content-Type")
	// res.Bodyの大きさ len(body) と同じ
	contentLength := res.ContentLength
	// リクエストURL
	reqURL := res.Request.URL.String()

	// GetParamResponse構造体初期化
	httpMethodsRes := HTTPMethodsResponse{}
	// GetParamResponse構造体に格納（JSON文字列をデコードし、結果を構造体へ）
	if err := json.Unmarshal(body, &httpMethodsRes); err != nil {
		log.Fatal(err)
	}
	// httpMethodsRes構造体に格納したものをJson形式へ
	jsonBytes, _ := json.Marshal(&httpMethodsRes)

	// レスポンスの表示（body は []byte バイト配列）
	fmt.Println("--TestAPIClientGetParam--Result-Start-------------------------------------------")
	fmt.Println("StatusCode：" + strconv.Itoa(statusCode)) // 数値を文字列に変換
	fmt.Println("Header Content-Type：" + header)
	fmt.Println("ContentLength：" + strconv.FormatInt(contentLength, 10)) // 数値int64を文字列に変換　引数設定10は10進数
	fmt.Println("RequestURL：" + reqURL)
	fmt.Println("Body：")
	fmt.Printf("%s", body)
	fmt.Println("--TestAPIClientGetParam--Result-End-------------------------------------------")

	// GetParamResponse構造体に格納した値を表示（デコードしたデータを表示）
	fmt.Println("--TestAPIClientGetParam--構造体-Start-------------------------------------------")
	// fmt.Println(httpMethodsRes.Args)
	fmt.Printf("code:%v, name:%v\n", httpMethodsRes.Args.Code, httpMethodsRes.Args.Name) // 構造体に格納されているので、取り出し何かしらの処理に利用できる
	fmt.Println("--TestAPIClientGetParam--構造体-End-------------------------------------------")
	// httpMethodsRes構造体に格納した値をJson形式にしたのを表示
	fmt.Println("--TestAPIClientGetParam--構造体-Json形式-Start-------------------------------------------")
	fmt.Println(string(jsonBytes))
	fmt.Println("--TestAPIClientGetParam--構造体-Json形式-End-------------------------------------------")

}

// TestAPIClientPost1 （パラメータがjson形式）
func TestAPIClientPost1() {
	// A simple HTTP Request & Response Service.「https://httpbin.org」を使用
	// 　URLに対してパラメータを追加したリクエストを送信すると、
	// 　クライアントのIPアドレスやユーザーエージェント、送信されたパラメータなどの情報がJSONデータとして返してくれる

	// httpのクライアントを作成する
	client := new(http.Client)

	// PostJSON構造体初期化
	postJSON := HTTPMethodsPost{}
	postJSON.ID = "1"
	postJSON.Name = "testTaro"
	// PostJSON構造をJson形式へ
	jsonBytes, _ := json.Marshal(&postJSON) // jsonBytesは[]byte型

	// リクエストを作成（Post）
	// 第一引数: http.MethodGet は "GET" を指定するのと同じ。メソッドを指定する（GET,POST,PUT,DELETE）
	// 第二引数: リクエストするURL
	// 第三引数: リクエストボディで フォーム、ファイル などのデータを送信する場合に使う。何も送信しない場合は nil を指定する
	req, err := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBuffer(jsonBytes)) // 「https://httpbin.org」サイトのHTTPMethodsのPOSTを使用
	if err != nil {
		panic(err.Error())
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

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
	header := res.Header.Get("Content-Type")
	// res.Bodyの大きさ len(body) と同じ
	contentLength := res.ContentLength
	// リクエストURL
	reqURL := res.Request.URL.String()

	// レスポンスの表示（body は []byte バイト配列）
	fmt.Println("--TestAPIClientPost1--Result-Start-------------------------------------------")
	fmt.Println("StatusCode：" + strconv.Itoa(statusCode)) // 数値を文字列に変換
	fmt.Println("Header Content-Type：" + header)
	fmt.Println("ContentLength：" + strconv.FormatInt(contentLength, 10)) // 数値int64を文字列に変換　引数設定10は10進数
	fmt.Println("RequestURL：" + reqURL)
	fmt.Println("Body：")
	fmt.Printf("%s", body)
	fmt.Println("--TestAPIClientPost1--Result-End-------------------------------------------")
}

// TestAPIClientPost2  （パラメータをURLエンコード）
func TestAPIClientPost2(postURL string, id string) {
	// A simple HTTP Request & Response Service.「https://httpbin.org」を使用
	// 　URLに対してパラメータを追加したリクエストを送信すると、
	// 　クライアントのIPアドレスやユーザーエージェント、送信されたパラメータなどの情報がJSONデータとして返してくれる

	// キー1=値1&キー2=値2&...を設定
	values := url.Values{} // "net/url"パッケージ
	values.Set("id", id)
	values.Add("token", "asdf1234")
	// fmt.Println(values.Encode()) // => id=99&token=asdf1234

	// httpのクライアントを作成する
	client := new(http.Client)
	// client := &http.Client{Timeout: time.Duration(30) * time.Second} // タイムアウトを30秒に指定してClient構造体を生成
	// client.Timeout = time.Duration(30) * time.Second // タイムアウトを30秒に指定

	// リクエストを作成（Post）
	// 第一引数: http.MethodGet は "GET" を指定するのと同じ。 メソッドを指定する（GET,POST,PUT,DELETE）
	// 第二引数: リクエストするURL
	// 第三引数: リクエストボディで フォーム、ファイル などのデータを送信する場合に使う。何も送信しない場合は nil を指定する
	req, err := http.NewRequest("POST", postURL, strings.NewReader(values.Encode())) // 「https://httpbin.org」サイトのHTTPMethodsのPOSTを使用
	if err != nil {
		panic(err.Error())
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
	header := res.Header.Get("Content-Type")
	// res.Bodyの大きさ len(body) と同じ
	contentLength := res.ContentLength
	// リクエストURL
	reqURL := res.Request.URL.String()

	// レスポンスの表示（body は []byte バイト配列）
	fmt.Println("--TestAPIClientPost2--Result-Start-------------------------------------------")
	fmt.Println("StatusCode：" + strconv.Itoa(statusCode)) // 数値を文字列に変換
	fmt.Println("Header Content-Type：" + header)
	fmt.Println("ContentLength：" + strconv.FormatInt(contentLength, 10)) // 数値int64を文字列に変換　引数設定10は10進数
	fmt.Println("RequestURL：" + reqURL)
	fmt.Println("Body：")
	fmt.Printf("%s", body)
	fmt.Println("--TestAPIClientPost2--Result-End-------------------------------------------")
}

// HTTPMethodsResponse 構造体（JSONデコード用に構造体定義）※「https://httpbin.org」サイトのHTTPMethodsのResponse bodyに合わすよう設定している
type HTTPMethodsResponse struct {
	Args    Args    //`json:"Args"`
	Headers Headers //`json:"Headers"`
	JSON    string  `json:"json"`
	Origin  string  `json:"origin"`
	URL     string  `json:"url"`
}

// Args  GetParamResponseのArgs
type Args struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Headers  GetParamResponseのHeaders
type Headers struct {
	AcceptEncoding string `json:"accept-encoding"`
	Host           string `json:"host"`
	UserAgent      string `json:"user-agent"`
	XAmznTraceID   string `json:"X-Amzn-Trace-Id"`
}

// HTTPMethodsPost  構造体　TestAPIClientPost1()（パラメータがjson形式）で使用
type HTTPMethodsPost struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
