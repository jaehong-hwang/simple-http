package database

import (
	"database/sql"

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
func (q *Query) Select(fields ...string) *Query {
	for _, field := range fields {
		q.selectors = append(q.selectors, field)
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

// Query run func
func (q *Query) Query(query string, args []interface{}) (*sql.Rows, error) {
	rows, err := q.connection.SQLDB.Query(query, args...)
	if err != nil {
		return nil,
			QueryError{
				QueryString: query,
				Parameters:  args,
				Message:     err.Error(),
			}
	}

	return rows, nil
}

// CRUD Querys
// ============

// Get rows by
// select query execute
func (q *Query) Get() (*sql.Rows, error) {
	query, args := command.
		NewSelect().
		Fields(q.selectors...).
		Where(q.where).
		ToString()

	return q.Query(query, args)
}

// Insert to table
func (q *Query) Insert() error {
	return nil
}

// Delete from table
func (q *Query) Delete() error {
	return nil
}
