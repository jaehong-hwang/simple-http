package database

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Pool struct
type Pool struct {
	Env          Env
	MaxOpenConns int
	MaxIdleConns int
	SQLDB        *sql.DB
}

// NewPool func
func NewPool(env Env, maxOpenConns, maxIdleConns int) (Pool, error) {
	db := Pool{
		Env:          env,
		MaxOpenConns: maxOpenConns,
		MaxIdleConns: maxIdleConns,
	}

	if err := db.Open(); err != nil {
		return Pool{}, err
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

	cols, _ := result.Rows.Columns()
	result.Rows.Next()

	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	if err := result.Rows.Scan(columnPointers...); err != nil {
		return err
	}

	values := make(map[string]interface{})
	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		values[colName] = *val
	}

	err = mapstructure.Decode(values, model)
	if err != nil {
		return err
	}

	return nil
}
