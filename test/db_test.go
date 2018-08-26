package test

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	sdb "github.com/jaehong-hwang/simple-http/database"
)

func connect(env sdb.Env) (*sql.DB, error) {
	db, err := sdb.Connect(env, 10, 10)
	return db, err
}

func TestConnection(t *testing.T) {
	env := GetDatabaseEnv()
	db, err := connect(env)
	if err != nil {
		t.Fatal(env.GetDataSourceName(), err.Error())
	}

	defer db.Close()

	t.Log("db connected")
}

type testStruct struct {
	f1 string
	f2 int
}

func TestQuery(t *testing.T) {
	env := GetDatabaseEnv()
	db, err := connect(env)
	if err != nil {
		t.Fatal(env.GetDataSourceName(), err.Error())
	}

	var f1 string
	var f2 int

	var datas []testStruct
	datas = append(datas, testStruct{f1: "a", f2: 1})
	datas = append(datas, testStruct{f1: "b", f2: 2})
	datas = append(datas, testStruct{f1: "c", f2: 3})
	_, err = db.Query("CREATE TABLE golang_query_test (f1 VARCHAR(20), f2 INT)")
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log("Insert Log ===")
	for _, data := range datas {
		t.Logf("f1 : %s, f2 : %d", data.f1, data.f2)
		_, err = db.Query("INSERT INTO golang_query_test (f1, f2) VALUES (?, ?)", data.f1, data.f2)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	rows, err := db.Query("SELECT * FROM golang_query_test")
	t.Log("Select result ===")
	for rows.Next() {
		rows.Scan(&f1, &f2)
		t.Logf("f1 : %s, f2 : %d", f1, f2)
	}

	_, err = db.Query("DROP TABLE golang_query_test")
	if err != nil {
		t.Fatal(err.Error())
	}
}
