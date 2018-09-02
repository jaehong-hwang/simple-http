package database

import (
	"database/sql"
	"time"
)

// QueryResult struct
type QueryResult struct {
	Rows        *sql.Rows
	QueryString string
	Parameters  []interface{}
	Duration    time.Duration
}
