package database

import (
	"database/sql"
	"strings"
)

// Query struct
type Query struct {
	connection *Pool
	table      string
	selectors  []string
	where      []map[string]interface{}
}

// From table setting
func (q *Query) From(table string) *Query {
	q.table = table
	return q
}

// Select append to model
func (q *Query) Select(args ...string) *Query {
	for _, arg := range args {
		q.selectors = append(q.selectors, arg)
	}

	return q
}

// Where append to model
func (q *Query) Where(query string, values ...interface{}) *Query {
	q.where = append(q.where, map[string]interface{}{"query": query, "args": values})

	return q
}

// CRUD Querys
// ============

// Get rows by
// select query execute
func (q *Query) Get() (*sql.Rows, error) {
	if len(q.selectors) < 1 {
		q.selectors = append(q.selectors, "*")
	}

	querystr := "SELECT `" + strings.Join(q.selectors, "`, `") + "` FROM `" + q.table + "`" + q.where

	return q.connection.SQLDB.Query(querystr)
}

// Insert to table
func (q *Query) Insert() error {
	return nil
}

// Delete from table
func (q *Query) Delete() error {
	return nil
}
