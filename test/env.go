package test

import (
	"os"

	shttp "github.com/jaehong-hwang/simple-http"
	sdb "github.com/jaehong-hwang/simple-http/database"
)

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
