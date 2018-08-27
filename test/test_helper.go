package test

import (
	"os"

	shttp "github.com/jaehong-hwang/simple-http"
	sdb "github.com/jaehong-hwang/simple-http/database"
)

// GetDBConnect func
func GetDBConnect(env sdb.Env) (sdb.Pool, error) {
	db, err := sdb.NewPool(env, 10, 10)
	return db, err
}

// GetServerEnv func
func GetServerEnv() shttp.ServerEnv {
	return shttp.ServerEnv{
		Port: 7857,
		Test: true,
	}
}

// GetDatabaseEnv Func
func GetDatabaseEnv() sdb.Env {
	return sdb.Env{
		Drive:    "mysql",
		Host:     "localhost",
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: "test",
		Charset:  "utf8mb4",
	}
}
