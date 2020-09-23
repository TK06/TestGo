package db

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 参考：https://qiita.com/smith-30/items/83ad73f1b746d6a12a47#%E6%8E%A5%E7%B6%9A%E4%B8%8A%E9%99%90%E3%81%AE%E8%A8%AD%E5%AE%9A

const slowQuery = "select sleep(5)"

func doQuery(db *gorm.DB) error {
	return db.Exec(slowQuery).Error
}

func Test_maxConn(t *testing.T) {

	type args struct {
		connCount int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				connCount: 100,
			},
		},
		{
			args: args{
				connCount: 101,
			},
		},
	}
	for _, tt := range tests {
		// DB接続をする
		InitDb()
		// DB接続を取得
		connPoolDB := ConnectionDb()
		// SetMaxIdleConnsはアイドル状態のコネクションプール内の最大数を設定します
		connPoolDB.DB().SetMaxIdleConns(0)
		// //SetMaxOpenConnsは接続済みのデータベースコネクションの最大数を設定
		// connPoolDB.DB().SetMaxOpenConns(1)

		if tt.args.connCount > 100 {
			connPoolDB.DB().SetMaxOpenConns(100)
		}

		t.Run(tt.name, func(t *testing.T) {
			db := connPoolDB
			wg := &sync.WaitGroup{}
			for index := 0; index < tt.args.connCount; index++ {
				go func() {
					wg.Add(1)
					defer wg.Done()
					if err := doQuery(db); err != nil {
						t.Errorf("%v\n", err)
					}
				}()
			}
			wg.Wait()
			connPoolDB.Close()
			// CloseDb() // db.goファイルのCloseDb()
		})
	}
}

// コネクションプールを使用しDBへの接続数の確認
func TestDb_useIdleConn(t *testing.T) {
	// DB接続をする
	InitDb()
	// DB接続を取得
	connPoolDB := ConnectionDb()
	// SetMaxIdleConnsはアイドル状態のコネクションプール内の最大数を設定します
	// connPoolDB.DB().SetMaxIdleConns(0)
	connPoolDB.DB().SetMaxIdleConns(5)
	// SetMaxOpenConnsは接続済みのデータベースコネクションの最大数を設定
	connPoolDB.DB().SetMaxOpenConns(100)

	sem := make(chan struct{}, 5)
	qs := getQueries(10)
	for _, item := range qs { // _, item := range qs は　インデックス, 値 := range qs　※ブランク識別子 _ を使ってインデックスを無視
		sem <- struct{}{}
		go func(db *gorm.DB, item string) {
			defer func() {
				<-sem
			}()
			if err := fetch(db, item); err != nil {
				panic(err)
			}
		}(connPoolDB, item)
	}
}

func getQueries(num int) []string {
	qs := make([]string, 0, num) // 組み込み関数make
	for index := 0; index < num; index++ {
		qs = append(qs, fmt.Sprintf("select sleep(%v)", random(0.001, 0.03)))
	}
	return qs
}
func random(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*(max-min) + min
}

// コネクションプールを使用したときのパフォーマンス比較
func BenchmarkUseIdleConn(b *testing.B) {
	// DB接続をする
	InitDb()
	// DB接続を取得
	connPoolDB := ConnectionDb()
	// SetMaxIdleConnsはアイドル状態のコネクションプール内の最大数を設定します
	// connPoolDB.DB().SetMaxIdleConns(0)
	connPoolDB.DB().SetMaxIdleConns(5)
	// SetMaxOpenConnsは接続済みのデータベースコネクションの最大数を設定
	connPoolDB.DB().SetMaxOpenConns(100)

	b.Run("", func(b *testing.B) {
		sem := make(chan struct{}, 5)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sem <- struct{}{}
			go func(db *gorm.DB) {
				defer func() {
					<-sem
				}()
				if err := fetch(db, "select sleep(0.01)"); err != nil {
					panic(err)
				}
			}(connPoolDB)
		}

	})
	connPoolDB.Close()
}

func fetch(db *gorm.DB, q string) error {
	return db.Exec(q).Error
}

// コネクションプールを使用したときのパフォーマンス比較
// 都度gorm.DBを作ってしまっている場合
func BenchmarkUseIdleConnForEveryTime(b *testing.B) {
	// DB接続をする
	InitDb()

	b.Run("", func(b *testing.B) {
		sem := make(chan struct{}, 5)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sem <- struct{}{}
			go func() {
				defer func() {
					<-sem
				}()
				// DB接続を取得
				connPoolDB := ConnectionDb()
				connPoolDB.DB().SetMaxOpenConns(100)
				connPoolDB.DB().SetMaxIdleConns(0)
				connPoolDB.DB().SetConnMaxLifetime(time.Hour)
				//defer connPoolDB.Close()
				if err := fetch(connPoolDB, "select sleep(0.01)"); err != nil {
					panic(err)
				}
			}()
		}
	})
}
