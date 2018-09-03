package command

import "fmt"

// Delete struc
type Delete struct {
	table string
	where *Where
}

// NewDelete func
func NewDelete(table string) *Delete {
	return &Delete{table: table}
}

// Where func
func (d *Delete) Where(where *Where) *Delete {
	d.where = where
	return d
}

// ToString func
func (d *Delete) ToString() (string, []interface{}) {
	var args []interface{}
	query := fmt.Sprintf("DELETE FROM `%s`", d.table)

	if d.where != nil {
		where, whereArgs := d.where.ToCommand()
		query += " " + where

		args = append(args, whereArgs...)
	}

	return query, args
}
