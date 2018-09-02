package database

import (
	"time"

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
func (q *Query) Query(query string, args []interface{}) (*QueryResult, error) {
	start := time.Now()
	rows, err := q.connection.SQLDB.Query(query, args...)
	elapsed := time.Since(start)

	result := &QueryResult{
		Rows:        rows,
		QueryString: query,
		Parameters:  args,
		Duration:    elapsed,
	}

	if err != nil {
		return result, err
	}

	return result, nil
}

// CRUD Querys
// ============

// Get rows by
// select query execute
func (q *Query) Get() (*QueryResult, error) {
	query, args := command.
		NewSelect(q.table).
		Fields(q.selectors...).
		Where(q.where).
		ToString()

	return q.Query(query, args)
}

// Insert to table
func (q *Query) Insert(rows ...map[string]interface{}) (*QueryResult, error) {
	insert := command.NewInsert(q.table)

	for i, row := range rows {
		if i == 0 {
			var fields []string
			for k := range row {
				fields = append(fields, k)
			}

			insert.SetFields(fields...)
		}

		var values []interface{}
		for _, val := range row {
			values = append(values, val)
		}

		insert.AddValues(values...)
	}

	query, args := insert.ToString()

	return q.Query(query, args)
}

// Delete from table
func (q *Query) Delete() error {
	return nil
}
