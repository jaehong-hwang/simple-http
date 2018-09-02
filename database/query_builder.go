package database

import (
	"database/sql"
	"strings"

	"github.com/jaehong-hwang/simple-http/database/command"
)

// Query struct
type Query struct {
	connection *Pool
	table      string
	selectors  []string
	where      *command.Where
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
	if q.where == nil {
		q.where = &command.Where{}
	}

	q.where = q.where.And(query, values...)

	return q
}

// OrWhere append to model
func (q *Query) OrWhere(query string, values ...interface{}) *Query {
	q.where = q.where.Or(query, values...)

	return q
}

// CRUD Querys
// ============

// Get rows by
// select query execute
func (q *Query) Get() (*sql.Rows, error) {
	var args []interface{}
	if len(q.selectors) < 1 {
		q.selectors = append(q.selectors, "*")
	}

	where, whereArgs := q.where.ToCommand()

	args = append(args, whereArgs...)

	querystr := "SELECT `" + strings.Join(q.selectors, "`, `") + "` FROM `" + q.table + "` " + where

	rows, err := q.connection.SQLDB.Query(querystr, args...)
	if err != nil {
		return nil,
			QueryError{
				QueryString: querystr,
				Parameters:  args,
				Message:     err.Error(),
			}
	}

	return rows, nil
}

// Insert to table
func (q *Query) Insert() error {
	return nil
}

// Delete from table
func (q *Query) Delete() error {
	return nil
}
