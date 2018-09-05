package command

import (
	"strconv"
	"strings"
)

// Select command in sql
type Select struct {
	table   string
	fields  []string
	where   *Where
	limit   []int
	orderBy []*OrderBy
}

// NewSelect struct return
func NewSelect(table string) *Select {
	return &Select{table: table}
}

// Fields set func
func (s *Select) Fields(fields ...string) *Select {
	s.fields = append(s.fields, fields...)
	return s
}

// Where func
func (s *Select) Where(w *Where) *Select {
	s.where = w
	return s
}

// Limit func
func (s *Select) Limit(limit ...int) *Select {
	s.limit = limit
	return s
}

// OrderBy func
func (s *Select) OrderBy(orderBy []*OrderBy) *Select {
	s.orderBy = orderBy
	return s
}

// ToString SELECT command
func (s *Select) ToString() (string, []interface{}) {
	var args []interface{}
	var fields string

	if len(s.fields) < 1 {
		fields = "*"
	} else {
		fields = "`" + strings.Join(s.fields, "`, `") + "`"
	}

	query := "SELECT " + fields + " FROM `" + s.table + "`"

	if s.where != nil {
		where, whereArgs := s.where.ToCommand()
		args = append(args, whereArgs...)
		query += " " + where
	}

	if len(s.orderBy) > 0 {
		var orders []string
		query += " ORDER BY "
		for _, order := range s.orderBy {
			orders = append(orders, order.Field+" "+order.Order)
		}
		query += strings.Join(orders, ", ")
	}

	if len(s.limit) > 0 && s.limit[0] != 0 {
		query += " LIMIT " + strconv.Itoa(s.limit[0])

		if len(s.limit) > 1 && s.limit[1] != 0 {
			query += ", " + strconv.Itoa(s.limit[1])
		}
	}

	return query, args
}
