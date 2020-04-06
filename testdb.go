package testdb

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"

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
func Setup(dbConf *DBConf, envPath string) error {
	var err error

	if dbConf == nil {
		dbConf, err = loadDBConf(envPath)
		if err != nil {
			return err
		}
	}

	if err = dbConf.checkValue(); err != nil {
		return err
	}

	DB, err = sql.Open(dbConf.DriverName, dbConf.createConf())
	if err != nil {
		return err
	}

	return DB.Ping()
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
