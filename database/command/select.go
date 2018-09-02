package command

import "strings"

// Select command in sql
type Select struct {
	table  string
	fields []string
	where  *Where
}

// NewSelect struct return
func NewSelect() *Select {
	s := new(Select)
	return s
}

// From set func
func (s *Select) From(table string) *Select {
	s.table = table
	return s
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

	query := "SELECT `" + strings.Join(s.fields, "`, `") + "` FROM `" + s.table + "`"

	if s.where != nil {
		where, whereArgs := s.where.ToCommand()
		args = append(args, whereArgs...)
		query += " " + where
	}

	return query, args
}
