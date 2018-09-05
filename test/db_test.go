package test

import (
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	sdb "github.com/jaehong-hwang/simple-http/database"
)

var env *sdb.Env
var db *sdb.Pool

func init() {
	var err error

	env = &sdb.Env{
		Drive:    "mysql",
		Host:     "localhost",
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: "test",
		Charset:  "utf8mb4",
	}

	db, err = sdb.NewPool(env, 10, 10)
	if err != nil {
		log.Fatal(env.GetDataSourceName(), err.Error())
	}
}

func TestConnection(t *testing.T) {
	err := db.Open()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	t.Log("db connected")
}

type testStruct struct {
	f1 string
	f2 int
}

func TestQuery(t *testing.T) {
	db.Open()
	defer db.Close()

	var f1 string
	var f2 int

	var datas []testStruct
	datas = append(datas, testStruct{f1: "a", f2: 1})
	datas = append(datas, testStruct{f1: "b", f2: 2})
	datas = append(datas, testStruct{f1: "c", f2: 3})
	_, err := db.SQLDB.Query("CREATE TABLE golang_query_test (f1 VARCHAR(20), f2 INT)")
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

func TestORMGet(t *testing.T) {
	db.Open()
	defer db.Close()

	result, err := db.Table("orm_test").
		Select("id", "f1", "f2").
		Where("id = ?", 1).
		OrWhere("id = ?", 2).
		Get()

	t.Logf("query string: \"%s\"", result.QueryString)
	t.Log("query parameters: ", result.Parameters)
	t.Log("query execution time: ", result.Duration)

	if err != nil {
		t.Fatalf("query error: %s", err.Error())
	}

	t.Log("==================================")

	var id int
	var f1, f2 string
	for result.Rows.Next() {
		result.Rows.Scan(&id, &f1, &f2)
		t.Logf("%2d | f1: %s | f2: %s", id, f1, f2)
	}
}

func TestORMInsert(t *testing.T) {
	db.Open()
	defer db.Close()

	result, err := db.Table("orm_test").Insert(
		map[string]interface{}{
			"f1": "test1",
			"f2": "test1",
		},
		map[string]interface{}{
			"f1": "test2",
			"f2": "test2",
		},
	)

	t.Logf("query string: \"%s\"", result.QueryString)
	t.Log("query parameters: ", result.Parameters)
	t.Log("query execution time: ", result.Duration)

	if err != nil {
		t.Fatalf("query error: %s", err.Error())
	}
}

func TestORMDelete(t *testing.T) {
	db.Open()
	defer db.Close()

	result, err := db.Table("orm_test").
		Where("f1 LIKE ?", "test%").
		Delete()

	t.Logf("query string: \"%s\"", result.QueryString)
	t.Log("query parameters: ", result.Parameters)
	t.Log("query execution time: ", result.Duration)

	if err != nil {
		t.Fatalf("query error: %s", err.Error())
	}
}

func TestORMUpdate(t *testing.T) {
	db.Open()
	defer db.Close()

	result, err := db.Table("orm_test").
		Where("f1 Like ?", "test%").
		Update(map[string]interface{}{
			"f1": "update1",
			"f2": "update1",
		})

	t.Logf("query string: \"%s\"", result.QueryString)
	t.Log("query parameters: ", result.Parameters)
	t.Log("query execution time: ", result.Duration)

	if err != nil {
		t.Fatalf("query error: %s", err.Error())
	}
}

// Board model
type Board struct {
	ID    int
	Title string
}

func TestModelGet(t *testing.T) {
	board := &Board{}
	err := db.GetByID(board, 1)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	t.Logf("query success")
	t.Logf("%d | %s", board.ID, board.Title)
}
