package database

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// Model struct
type Model struct {
	table string
	V     interface{}
}

var pool *Pool

// UsePool for model
func UsePool(p *Pool) {
	pool = p
}

// NewModel func
func NewModel(name string, vo interface{}) *Model {
	m := &Model{table: name, V: vo}
	return m
}

// Get by id
func (m *Model) Get(id int) (*Model, error) {
	err := pool.Open()
	if err != nil {
		return nil, err
	}

	defer pool.Close()

	result, err := pool.Table(m.table).Where("id=?", id).Get()
	if err != nil {
		return nil, err
	}

	cols, _ := result.Rows.Columns()
	result.Rows.Next()

	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	if err := result.Rows.Scan(columnPointers...); err != nil {
		return nil, err
	}

	values := make(map[string]interface{})
	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		values[colName] = *val
	}

	err = mapstructure.Decode(values, m.V)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}
