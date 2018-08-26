package test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	sdb "github.com/jaehong-hwang/simple-http/database"
)

var db sdb.Pool

func TestConnection(t *testing.T) {
	var err error
	env := GetDatabaseEnv()

	if db, err = sdb.NewPool(env, 10, 10); err != nil {
		t.Fatal(env.GetDataSourceName(), err.Error())
	}

	defer db.Close()

	t.Log("db connected")
}
