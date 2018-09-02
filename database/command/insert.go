package command

import (
	"fmt"
	"strings"
)

// Insert command struct
type Insert struct {
	table  string
	fields []string
	values []interface{}
}

// NewInsert func
func NewInsert(table string) *Insert {
	return &Insert{table: table}
}

// SetFields func
func (i *Insert) SetFields(fields ...string) *Insert {
	i.fields = fields
	return i
}

// AddValues func
func (i *Insert) AddValues(values ...interface{}) *Insert {
	i.values = append(i.values, values...)
	return i
}

// ToString func
func (i *Insert) ToString() (string, []interface{}) {
	query := fmt.Sprintf("INSERT INTO `%s` (`"+strings.Join(i.fields, "`, `")+"`) VALUES ", i.table)

	var argval []string
	for j := len(i.fields); j > 0; j-- {
		argval = append(argval, "?")
	}

	var values []string
	for j := len(i.values) / len(i.fields); j > 0; j-- {
		values = append(values, "("+strings.Join(argval, ", ")+")")
	}

	query += strings.Join(values, ", ")

	return query, i.values
}
