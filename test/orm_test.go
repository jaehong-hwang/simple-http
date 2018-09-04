package test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jaehong-hwang/simple-http/database"
)

// Board model
type Board struct {
	*database.Model
}

func TestORMGet(t *testing.T) {
	env := GetDatabaseEnv()
	db, err := GetDBConnect(env)
	if err != nil {
		t.Fatal(env.GetDataSourceName(), err.Error())
	}

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
	env := GetDatabaseEnv()
	db, err := GetDBConnect(env)
	if err != nil {
		t.Fatal(env.GetDataSourceName(), err.Error())
	}

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
	env := GetDatabaseEnv()
	db, err := GetDBConnect(env)
	if err != nil {
		t.Fatal(env.GetDataSourceName(), err.Error())
	}

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
	env := GetDatabaseEnv()
	db, err := GetDBConnect(env)
	if err != nil {
		t.Fatal(env.GetDataSourceName(), err.Error())
	}

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

func TestModelGet(t *testing.T) {
	env := GetDatabaseEnv()
	_, err := GetDBConnect(env)
	if err != nil {
		t.Fatal(env.GetDataSourceName(), err.Error())
	}

	type BoardVO struct {
		ID    int
		Title []uint8
	}

	Board := database.NewModel("board", &BoardVO{})
	result, err := Board.Get(1)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	values := result.V.(*BoardVO)

	t.Log(values.ID, values.Title)
}
