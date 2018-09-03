package command

import (
	"fmt"
	"strings"
)

// Update struct
type Update struct {
	table string
	where *Where
	set   *map[string]interface{}
}

// NewUpdate func
func NewUpdate(table string) *Update {
	return &Update{table: table}
}

// Set func
func (u *Update) Set(set *map[string]interface{}) *Update {
	u.set = set
	return u
}

// Where func
func (u *Update) Where(where *Where) *Update {
	u.where = where
	return u
}

// ToString func
func (u *Update) ToString() (string, []interface{}) {
	var args []interface{}
	query := fmt.Sprintf("UPDATE `%s` SET ", u.table)

	var set []string
	for key, val := range *u.set {
		set = append(set, fmt.Sprintf("%s=?", key))
		args = append(args, val)
	}

	query += strings.Join(set, ", ")

	if u.where != nil {
		where, whereArgs := u.where.ToCommand()

		query += " " + where
		args = append(args, whereArgs...)
	}

	return query, args
}
