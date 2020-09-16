package api

import (
	"TestGo/db"

	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin" // フレームワーク Gin
	"github.com/jinzhu/gorm"   // GORM Go言語のORマッパー
)

// API CRUD(GET,POST,PUT,DELETE)
func API() {

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

		// DB接続を取得
		db := db.ConnectionDb()

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

		// DB接続を取得
		db := db.ConnectionDb()

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

		// DB接続を取得
		db := db.ConnectionDb()

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

// User テーブル　structを用意する。
type User struct {
	gorm.Model
	Name string `json:"name"`
	Age  int    `json:"age"`
}
