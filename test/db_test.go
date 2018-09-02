package test

import (
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestConnection(t *testing.T) {
	env := GetDatabaseEnv()
	db, err := GetDBConnect(env)
	if err != nil {
		t.Fatal(env.GetDataSourceName(), err.Error())
	}

	db.Open()
	defer db.Close()

	t.Log("db connected")

	time.Sleep(3 * time.Second)
}

type testStruct struct {
	f1 string
	f2 int
}

func TestQuery(t *testing.T) {
	env := GetDatabaseEnv()
	db, err := GetDBConnect(env)
	if err != nil {
		t.Fatal(env.GetDataSourceName(), err.Error())
	}

	db.Open()
	defer db.Close()

	var f1 string
	var f2 int

	var datas []testStruct
	datas = append(datas, testStruct{f1: "a", f2: 1})
	datas = append(datas, testStruct{f1: "b", f2: 2})
	datas = append(datas, testStruct{f1: "c", f2: 3})
	_, err = db.SQLDB.Query("CREATE TABLE golang_query_test (f1 VARCHAR(20), f2 INT)")
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log("Insert Log ===")
	for _, data := range datas {
		t.Logf("f1 : %s, f2 : %d", data.f1, data.f2)
		_, err = db.SQLDB.Query("INSERT INTO golang_query_test (f1, f2) VALUES (?, ?)", data.f1, data.f2)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	rows, err := db.SQLDB.Query("SELECT * FROM golang_query_test")
	t.Log("Select result ===")
	for rows.Next() {
		rows.Scan(&f1, &f2)
		t.Logf("f1 : %s, f2 : %d", f1, f2)
	}

	_, err = db.SQLDB.Query("DROP TABLE golang_query_test")
	if err != nil {
		t.Fatal(err.Error())
	}
}
