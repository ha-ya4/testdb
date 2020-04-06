package testdb

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"time"

	"github.com/joho/godotenv"
)

var (
	// DB データベース
	DB *sql.DB
	// TablesName 削除したいデータが入ってるテーブル名のスライス
	TablesName []string
)

var (
	// ErrNoTable テーブルのデータを削除する関数の引数とグローバル変数TablesName両方nilのときに返すエラー
	ErrNoTable = errors.New("missing table name")
)

// Setup 指定したenvファイルからDBの接続情報を取得してDBに接続しグローバル変数DBにセットする
func Setup(dbConf *DBConf) error {
	var err error

	if err = dbConf.checkValue(); err != nil {
		return err
	}

	DB, err = sql.Open(dbConf.DriverName, dbConf.createConf())
	if err != nil {
		return err
	}

	return DB.Ping()
}

// SetupByEnv 指定したenvファイルからDBの接続情報を取得してDBに接続しグローバル変数DBにセットする
func SetupByEnv(envPath string) error {
	dbConf, err := loadDBConf(envPath)
	if err != nil {
		return err
	}
	return Setup(dbConf)
}

// DBConf DBに接続するための設定
type DBConf struct {
	DriverName string
	UserName   string
	Password   string
	DBName     string
}

func loadDBConf(envPath string) (*DBConf, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		return nil, err
	}

	c := &DBConf{
		DriverName: os.Getenv("DB_DRIVER_NAME"),
		UserName:   os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASS"),
		DBName:     os.Getenv("DB_NAME"),
	}

	return c, err
}

func (c *DBConf) checkValue() error {
	value := reflect.Indirect(reflect.ValueOf(c))
	typ := reflect.TypeOf(*c)

	var fn []string
	for i := 0; i < typ.NumField(); i++ {
		// i番目のフィールドの値が初期値ならそのフィールド名をスライスに入れる
		if value.Field(i).IsZero() {
			fn = append(fn, typ.Field(i).Name)
		}
	}

	var err error
	if len(fn) != 0 {
		e := "Err: db conf "
		for _, n := range fn {
			e = e + n + " "
		}
		err = errors.New(e + "fields have no value")
	}

	return err
}

func (c *DBConf) createConf() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s", c.UserName, c.Password, c.DBName)

}

// DeleteTablesFrom 引数で受け取ったテーブル名のテーブルのデータを全て削除する。
func DeleteTablesFrom(tablesName []string) error {
	var err error

	for _, n := range tablesName {
		_, err := DB.Exec(fmt.Sprintf("DELETE FROM %s", n))
		if err != nil {
			return err
		}
	}

	return err
}

// DeleteTables グローバル変数にセットされているテーブル名のテーブルデータを全て削除する
func DeleteTables() error {
	if len(TablesName) == 0 {
		return ErrNoTable
	}
	return DeleteTablesFrom(TablesName)
}

// DeleteFrom 引数で受け取ったテーブル名のテーブルからデータを全て削除する
func DeleteFrom(tableName string) (sql.Result, error) {
	return DB.Exec(fmt.Sprintf("DELETE FROM %s", tableName))
}

// RandNum 引数に指定した上限の範囲内でランダムな整数を返す
func RandNum(limit int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(limit)
}

// RandstrFrom 文字数を引数lengthとして引数cの文字列を使ったランダムな文字列にして返す
func RandstrFrom(length int, c string) string {
	rand.Seed(time.Now().UnixNano())

	bytes := make([]byte, length)
	for i := range bytes {
		n := rand.Intn(len(c))
		bytes[i] = c[n]
	}
	return string(bytes)
}

const character = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Randstr 文字数を引数lenとしてabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZをランダムに組み合わせた文字列を返す
func Randstr(length int) string {
	return RandstrFrom(length, character)
}
