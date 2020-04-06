package testdb

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func insert111() (sql.Result, error) {
	return DB.Exec("INSERT INTO testdb_schema.test_table111 VALUES ('testtesttest')")
}

func insert222() (sql.Result, error) {
	return DB.Exec("INSERT INTO testdb_schema.test_table222 VALUES ('testtesttest')")
}

func select111() error {
	row := DB.QueryRow("SELECT * FROM testdb_schema.test_table111")
	var col string
	return row.Scan(&col)
}

func select222() error {
	row := DB.QueryRow("SELECT * FROM testdb_schema.test_table222")
	var col string
	return row.Scan(&col)
}

func TestMain(m *testing.M) {
	err := SetupByEnv("./test.env")
	if err != nil {
		fmt.Println(err)
		return
	}

	c := m.Run()

	os.Exit(c)
}

func TestSetup(t *testing.T) {
	t.Skip()

	c := &DBConf{
		DriverName: "postgres",
		UserName:   "", // テストに使う設定を入れる
		Password:   "", //
		DBName:     "testdb", //
	}
	err := Setup(c)
	assert.NoError(t, err)

	c.UserName = ""
	err = Setup(c)
	assert.Error(t, err)
}

func TestSetupByEnv(t *testing.T) {
	t.Skip()

	err := SetupByEnv("./test.env")
	assert.NoError(t, err)
}

func TestDeleteTablesFrom(t *testing.T) {
	_, err := insert111()
	_, err = insert222()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	// insetできてるか一応確認
	err = select111()
	assert.NoError(t, err, sql.ErrNoRows)
	err = select222()
	assert.NoError(t, err, sql.ErrNoRows)

	tablesName := []string{
		"testdb_schema.test_table111",
		"testdb_schema.test_table222",
	}

	err = DeleteTablesFrom(tablesName)
	assert.NoError(t, err)

	// データが削除されてsql: no rows in result setが返ってくるか
	err = select111()
	assert.Exactly(t, err, sql.ErrNoRows)
	err = select222()
	assert.Exactly(t, err, sql.ErrNoRows)
}

func TestDeleteTables(t *testing.T) {
	_, err := insert111()
	_, err = insert222()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	// insetできてるか一応確認
	err = select111()
	assert.NoError(t, err, sql.ErrNoRows)
	err = select222()
	assert.NoError(t, err, sql.ErrNoRows)

	// グローバル変数にセット
	TablesName = []string{
		"testdb_schema.test_table111",
		"testdb_schema.test_table222",
	}

	err = DeleteTables()
	assert.NoError(t, err)

	// データが削除されてsql: no rows in result setが返ってくるか
	err = select111()
	assert.Exactly(t, err, sql.ErrNoRows)
	err = select222()
	assert.Exactly(t, err, sql.ErrNoRows)
}

func TestDeleteFrom(t *testing.T) {
	_, err := insert111()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	// insetできてるか一応確認
	err = select111()
	assert.NoError(t, err, sql.ErrNoRows)

	tableName := "testdb_schema.test_table111"

	_, err = DeleteFrom(tableName)
	assert.NoError(t, err)

	// データが削除されてsql: no rows in result setが返ってくるか
	err = select111()
	assert.Exactly(t, err, sql.ErrNoRows)
}

func TestDBConfCheckValue(t *testing.T) {
	c := &DBConf{
		DriverName: "db",
		UserName:   "testuser",
		Password:   "testpass",
		DBName:     "testdb",
	}

	err := c.checkValue()
	assert.NoError(t, err)

	c.DriverName = ""
	c.Password = ""
	err = c.checkValue()
	expect := errors.New("Err: db conf DriverName Password fields have no value")
	assert.Exactly(t, err, expect)
}
