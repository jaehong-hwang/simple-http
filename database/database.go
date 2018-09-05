package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// Pool struct
type Pool struct {
	Env          *Env
	MaxOpenConns int
	MaxIdleConns int
	SQLDB        *sql.DB
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
	return p.SQLDB.Close()
}

// Table func
// query start with table name
func (p *Pool) Table(table string) *Query {
	query := &Query{connection: p}
	return query.From(table)
}

// Query start
func (p *Pool) Query() *Query {
	return &Query{connection: p}
}

// ORM Methods
// ==============

// Get by id
func (p *Pool) Get(model interface{}, id int) error {
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
