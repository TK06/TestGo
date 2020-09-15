package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin" // フレームワーク Gin
	"github.com/jinzhu/gorm"   // GORM Go言語のORマッパー
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	// APIクライアントのテスト
	// A simple HTTP Request & Response Service.「https://httpbin.org」を使用
	// 　URLに対してパラメータを追加したリクエストを送信すると、
	// 　クライアントのIPアドレスやユーザーエージェント、送信されたパラメータなどの情報がJSONデータとして返してくれる
	// 単純なGET
	testAPIClientGet()
	// パラメーター付きGET
	//testAPIClientGetParam()
	// POST
	//testAPIClientPost()

	// DB接続（MySQL）
	db := getGormConnect()
	// 実行完了後DB接続を閉じる
	defer db.Close()

	// gin の変数を定義しています
	// デフォルトのミドルウェアとともにginルーターを作成
	// Logger と アプリケーションクラッシュをキャッチするRecoveryミドルウェア を保有しています
	router := gin.Default()

	// CREATE
	// router.POST("/post", func(c *gin.Context) { //確認　curl -X POST -H "Content-Type: application/json" localhost:8080/post
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "POST",
	// 	})
	// })
	// POSTを切り出して呼び出し
	router.POST("/post", post) //確認　curl -X POST -H "Content-Type: application/json" -d "{\"id\":\"1\", \"name\":\"hoge\"}" localhost:8080/post

	// CREATE DB登録
	router.POST("/post/user", func(ctx *gin.Context) { //確認　curl -X POST -H "Content-Type: application/json" -d "{\"name\":\"test0\", \"age\":10}" localhost:8080/post/user
		// user構造体の初期化
		user := User{}
		now := time.Now()
		user.CreatedAt = now
		user.UpdatedAt = now

		err := ctx.BindJSON(&user)
		if err != nil {
			ctx.String(http.StatusBadRequest, "Request is failed: "+err.Error())
		}

		db.NewRecord(user)
		// 登録
		db.Create(&user) // 引数でアドレスを指定しないといけないみたい
		//db.Create(user) // レコードを登録できない。「using unaddressable value」とコンソールに出力されている。ナニコレ
		if db.NewRecord(user) == false {
			ctx.JSON(http.StatusOK, user)
		}
	})

	// READ
	router.GET("/get", func(ctx *gin.Context) { // curl -X GET localhost:8080/get
		ctx.JSON(http.StatusOK, gin.H{
			"message": "GET",
		})
	})
	// // READ パラメーター付き (/get/XXX)
	// router.GET("/get/:number", func(ctx *gin.Context) {
	// 	number := ctx.Param("number")
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"message": "GETparam",
	// 		"number":  number,
	// 	})
	// })
	// READ DB 1レコード取得
	router.GET("/get/user/:id", func(ctx *gin.Context) { //確認　curl -X GET localhost:8080/get/user/1
		user := User{}
		id := ctx.Param("id")

		db.Where("id = ?", id).Find(&user)
		ctx.JSON(http.StatusOK, user)
	})

	// UPDATE
	router.PUT("/put/user/:id", func(ctx *gin.Context) {
		// 確認　curl -X PUT localhost:8080/put/user/1 →ctx.String(http.StatusBadRequest, "Request is failed: "+err.Error())　に入る。HTTPステータス400
		// 確認　curl -X PUT -H "Content-Type: application/json" -d "{\"name\":\"\", \"age\":21}" localhost:8080/put/user/1　→HTTPステータス200

		// user構造体の初期化
		user := User{}
		// パスパラメータから取得
		id := ctx.Param("id")

		if err := ctx.BindJSON(&user); err != nil {
			ctx.String(http.StatusBadRequest, "Request is failed: "+err.Error()) // jsonリクエストパラメータがない場合ここにくる
		}
		// Jsonリクエストパラメータがuser構造体に設定されたので、取得する
		updateName := user.Name
		updateAge := user.Age
		fmt.Println("Name: " + updateName)
		fmt.Println("Age:", updateAge)

		// テーブルからレコード取得　（取得したものがuserに入る）
		// db.Where("id = ?", id).First(&user) // ORDER BY id ASC LIMIT 1　が付いている
		// fmt.Println(user)
		db.Where("id = ?", id).Find(&user) // ORDER BY id ASC LIMIT 1　が付いていない
		fmt.Println(user)

		// userにJsonリクエストパラメータから取得した値を設定する
		user.Name = updateName
		user.Age = updateAge
		// 更新（update set ）
		db.Save(&user) // すべてのフィールドが保存される。UPDATE user SET updated_at =, name = '●●', age = ●● WHERE id = '●●'　（GORMでupdated_atは自動で更新してくれる）
		// db.Updates(&user) // 更新できない。コンソールに「Error 1103: Incorrect table name '' 」テーブル名が無いって怒られている。

		// レスポンス　HTTPステータスコード200　Jsonでupdateしたuserテーブルのレコード情報を返す
		ctx.JSON(http.StatusOK, user)
	})

	// DELETE
	router.DELETE("/delete", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "DELETE",
		})
	})

	// PORT環境変数が指定されていない場合は8080ポートで待受
	router.Run()
	// router.Run(":8080") // ポートをハードコーディングした場合
}

// RequestJSON CREATEのPOSTで使用
type RequestJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CREATEのPOSTで使用
func post(ctx *gin.Context) {
	var json RequestJSON

	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": json.ID, "name": json.Name})
}

// DB接続（MySQL）
func getGormConnect() *gorm.DB {
	DBMS := "mysql"
	USER := "root"                    // DB User
	PASS := "root"                    // DB Password
	PROTOCOL := "tcp(127.0.0.1:3306)" // tcp(DBのIPアドレス:Port)
	DBNAME := "world"                 // DB名
	// GoでMySQLのtimestampを読みこむ方法：データベース接続時に「?parseTime=true」を指定する
	// parseTime=true を入れておかないと、time.Time型のメンバに値を取り込む際にエラーとなってしまう。
	// またデフォルトではタイムゾーンが世界標準時となるので、指定したいのなら loc=Asia/Tokyo のようにする(URIエスケープが必要)。
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?parseTime=true&loc=Asia%2FTokyo"
	db, err := gorm.Open(DBMS, CONNECT)

	// エラーハンドリング
	if err != nil {
		panic(err.Error())
	}
	// DBエンジンを「InnoDB」に設定
	db.Set("gorm:table_options", "ENGINE=InnoDB")
	// ログ出力を有効
	db.LogMode(true)
	// 登録するテーブル名を単数形にする（デフォルトは複数形）
	// ※このmain.goファイルでいくと、User構造体(userテーブルがマイグレーション(db.AutoMigrate(&User{}))で作られる。単数形の設定がないと複数形の名前でusersテーブルが作成される)
	db.SingularTable(true)
	// マイグレーション（テーブルが無い時は自動生成）
	db.AutoMigrate(&User{})

	fmt.Println("db connected: ", &db)
	return db
}

// User テーブル　structを用意する。
type User struct {
	gorm.Model
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// testAPIClientGet 　net/httpパッケージを使用
func testAPIClientGet() {

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
