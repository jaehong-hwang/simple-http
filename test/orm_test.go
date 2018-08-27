package test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestORMGet(t *testing.T) {
	env := GetDatabaseEnv()
	db, err := GetDBConnect(env)
	if err != nil {
		t.Fatal(env.GetDataSourceName(), err.Error())
	}

	rows, err := db.Query().
		From("orm_test").
		Select("id", "f1", "f2").
		Where("id = ?", 1).
		Get()

	if err != nil {
		t.Fatalf("query error: %s", err.Error())
	}

	var id int
	var f1, f2 string
	for rows.Next() {
		rows.Scan(&id, &f1, &f2)
		t.Logf("%2d | f1: %s | f2: %s", id, f1, f2)
	}
}
