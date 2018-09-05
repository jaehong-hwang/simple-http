package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jaehong-hwang/simple-http/database/command"
)

// Pool struct
type Pool struct {
	Env          *Env
	MaxOpenConns int
	MaxIdleConns int
	SQLDB        *sql.DB

	table     string
	selectors []string
	where     *command.Where
}

// NewPool func
func NewPool(env *Env, maxOpenConns, maxIdleConns int) (*Pool, error) {
	db := &Pool{
		Env:          env,
		MaxOpenConns: maxOpenConns,
		MaxIdleConns: maxIdleConns,
	}

	if err := db.Open(); err != nil {
		return &Pool{}, err
	}

	db.Close()

	return db, nil
}

// Open DB Pool
func (p *Pool) Open() error {
	var err error
	p.SQLDB, err = sql.Open(p.Env.Drive, p.Env.GetDataSourceName())

	if err != nil {
		return err
	}

	if err = p.SQLDB.Ping(); err != nil {
		return err
	}

	p.SQLDB.SetMaxOpenConns(p.MaxOpenConns)
	p.SQLDB.SetMaxIdleConns(p.MaxIdleConns)
	return nil
}

// Close pool
func (p *Pool) Close() error {
	p.table = ""
	p.selectors = []string{}
	p.where = nil

	return p.SQLDB.Close()
}

// Table func
// query start with table name
func (p *Pool) Table(table string) *Pool {
	p.table = table
	return p
}

// Select append to model
func (p *Pool) Select(fields ...string) *Pool {
	for _, field := range fields {
		p.selectors = append(p.selectors, field)
	}

	return p
}

// Where append to model
func (p *Pool) Where(query string, values ...interface{}) *Pool {
	if p.where == nil {
		p.where = &command.Where{}
	}

	p.where = p.where.And(query, values...)

	return p
}

// OrWhere append to model
func (p *Pool) OrWhere(query string, values ...interface{}) *Pool {
	p.where = p.where.Or(query, values...)

	return p
}

// Query run func
func (p *Pool) Query(query string, args []interface{}) (*QueryResult, error) {
	start := time.Now()
	rows, err := p.SQLDB.Query(query, args...)
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
func (p *Pool) Get() (*QueryResult, error) {
	query, args := command.
		NewSelect(p.table).
		Fields(p.selectors...).
		Where(p.where).
		ToString()

	return p.Query(query, args)
}

// Insert to table
func (p *Pool) Insert(rows ...map[string]interface{}) (*QueryResult, error) {
	insert := command.NewInsert(p.table)

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

	return p.Query(query, args)
}

// Update to table
func (p *Pool) Update(val map[string]interface{}) (*QueryResult, error) {
	query, args := command.NewUpdate(p.table).
		Where(p.where).
		Set(&val).
		ToString()

	return p.Query(query, args)
}

// Delete from table
func (p *Pool) Delete() (*QueryResult, error) {
	query, args := command.
		NewDelete(p.table).
		Where(p.where).
		ToString()

	return p.Query(query, args)
}

// ORM Methods
// ==============

// GetByID func
func (p *Pool) GetByID(model interface{}, id int) error {
	err := p.Open()
	if err != nil {
		return err
	}

	defer p.Close()

	names := strings.Split(reflect.ValueOf(model).Type().String(), ".")
	result, err := p.Table(names[len(names)-1]).Where("id=?", id).Get()
	if err != nil {
		return err
	}

	defer result.Rows.Close()

	result.Rows.Next()

	types, err := result.Rows.ColumnTypes()
	if err != nil {
		return err
	}

	modelValue := reflect.ValueOf(model)

	ptrs := make([]interface{}, len(types))
	for i, tp := range types {
		field := modelValue.Elem().FieldByName(toPascalCase(tp.Name()))
		if !field.IsValid() {
			return fmt.Errorf("Interface `%s` does not have the field `%s`", modelValue.Type(), tp.Name())
		}

		ptrs[i] = field.Addr().Interface()
	}

	err = result.Rows.Scan(ptrs...)
	if err != nil {
		return err
	}

	return nil
}

func toPascalCase(str string) string {
	if str == "id" {
		return "ID"
	}

	return strings.ToUpper(string(str[0])) + string(str[1:])
}
