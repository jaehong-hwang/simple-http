package database

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jaehong-hwang/simple-http/database/command"
)

// Connection of DB Pool
type Connection struct {
	conn    *sql.Conn
	context *context.Context
	table   string
	where   *command.Where

	//select fields
	selectors []string
	orderBy   []*command.OrderBy
	limit     int
	offset    int
}

// Close connection
func (c *Connection) Close() error {
	return c.conn.Close()
}

// Table func
// query start with table name
func (c *Connection) Table(table string) *Connection {
	c.table = table
	return c
}

// Select append to model
func (c *Connection) Select(fields ...string) *Connection {
	for _, field := range fields {
		c.selectors = append(c.selectors, field)
	}

	return c
}

// Where append to model
func (c *Connection) Where(query string, values ...interface{}) *Connection {
	if c.where == nil {
		c.where = &command.Where{}
	}

	c.where = c.where.And(query, values...)

	return c
}

// OrWhere append to model
func (c *Connection) OrWhere(query string, values ...interface{}) *Connection {
	c.where = c.where.Or(query, values...)

	return c
}

// Limit func
func (c *Connection) Limit(limit ...int) *Connection {
	if len(limit) >= 2 {
		c.offset = limit[1]
	}

	if len(limit) >= 1 {
		c.limit = limit[0]
	}

	return c
}

// OrderBy Func
func (c *Connection) OrderBy(field string, order string) *Connection {
	c.orderBy = append(c.orderBy, &command.OrderBy{Field: field, Order: order})
	return c
}

// Query run func
func (c *Connection) Query(query string, args []interface{}) (*QueryResult, error) {
	start := time.Now()
	rows, err := c.conn.QueryContext(*c.context, query, args...)
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
func (c *Connection) Get() (*QueryResult, error) {
	query, args := command.
		NewSelect(c.table).
		Fields(c.selectors...).
		Where(c.where).
		Limit(c.limit, c.offset).
		OrderBy(c.orderBy).
		ToString()

	return c.Query(query, args)
}

// Insert to table
func (c *Connection) Insert(rows ...map[string]interface{}) (*QueryResult, error) {
	insert := command.NewInsert(c.table)

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

	result, err := c.Query(query, args)
	defer result.Rows.Close()

	return result, err
}

// Update to table
func (c *Connection) Update(val map[string]interface{}) (*QueryResult, error) {
	query, args := command.NewUpdate(c.table).
		Where(c.where).
		Set(&val).
		ToString()

	result, err := c.Query(query, args)
	defer result.Rows.Close()

	return result, err
}

// Delete from table
func (c *Connection) Delete() (*QueryResult, error) {
	query, args := command.
		NewDelete(c.table).
		Where(c.where).
		ToString()

	result, err := c.Query(query, args)
	defer result.Rows.Close()

	return result, err
}

// ORM Methods
// ==============

// GetByID func
func (c *Connection) GetByID(model interface{}, id int) error {
	names := strings.Split(reflect.ValueOf(model).Type().String(), ".")
	result, err := c.Table(names[len(names)-1]).Where("id=?", id).Get()
	if err != nil {
		return err
	}

	defer result.Rows.Close()

	for result.Rows.Next() {
		err = scanModel(model, result.Rows)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("%s model have no item by id %d", names, id)
}

// First func
func (c *Connection) First(model interface{}) error {
	names := strings.Split(reflect.ValueOf(model).Type().String(), ".")
	result, err := c.Table(names[len(names)-1]).Limit(1).Get()
	if err != nil {
		return err
	}

	defer result.Rows.Close()

	result.Rows.Next()

	err = scanModel(model, result.Rows)
	if err != nil {
		return err
	}

	return nil
}

// get model by pointers
func scanModel(model interface{}, rows *sql.Rows) error {
	modelValue := reflect.ValueOf(model)

	types, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	ptrs := make([]interface{}, len(types))
	for i, tp := range types {
		field := modelValue.Elem().FieldByName(toPascalCase(tp.Name()))
		if !field.IsValid() {
			return fmt.Errorf("Interface `%s` does not have the field `%s`", modelValue.Type(), tp.Name())
		}

		ptrs[i] = field.Addr().Interface()
	}

	err = rows.Scan(ptrs...)
	if err != nil {
		return err
	}

	return nil
}

// string to pascal case
func toPascalCase(str string) string {
	if str == "id" {
		return "ID"
	}

	return strings.ToUpper(string(str[0])) + string(str[1:])
}
