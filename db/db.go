package db

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"                  // GORM Go言語のORマッパー
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysqlドライバ
)

var db *gorm.DB
var err error

// InitDb DB接続をする（MySQL）
// func GetGormConnect() *gorm.DB {
func InitDb() {
	DBMS := "mysql"
	USER := "root"                    // DB User
	PASS := "root"                    // DB Password
	PROTOCOL := "tcp(127.0.0.1:3306)" // tcp(DBのIPアドレス:Port)
	DBNAME := "world"                 // DB名
	// GoでMySQLのtimestampを読みこむ方法：データベース接続時に「?parseTime=true」を指定する
	// parseTime=true を入れておかないと、time.Time型のメンバに値を取り込む際にエラーとなってしまう。
	// またデフォルトではタイムゾーンが世界標準時となるので、指定したいのなら loc=Asia/Tokyo のようにする(URIエスケープが必要)。
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?parseTime=true&loc=Asia%2FTokyo"
	db, err = gorm.Open(DBMS, CONNECT)

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
	// // マイグレーション（テーブルが無い時は自動生成）
	// db.AutoMigrate(&User{})

	// コネクションプール設定
	// SetMaxIdleConnsはアイドル状態のコネクションプール内の最大数を設定します
	db.DB().SetMaxIdleConns(10)
	// SetMaxOpenConnsは接続済みのデータベースコネクションの最大数を設定します
	db.DB().SetMaxOpenConns(100)
	// SetConnMaxLifetimeは再利用され得る最長時間を設定します
	db.DB().SetConnMaxLifetime(time.Hour)

	fmt.Println("db connected: ", &db)
	fmt.Println("DB接続成功")
	// return db
}

// ConnectionDb DB接続を取得
func ConnectionDb() *gorm.DB {
	return db
}

// CloseDb DB切断
func CloseDb() {
	if db != nil {
		db.Close()

	}
}
