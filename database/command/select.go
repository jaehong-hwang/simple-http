package command

import "strings"

// Select command in sql
type Select struct {
	table  string
	fields []string
	where  *Where
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

	return query, args
}
