package command

// Where command struct
type Where struct {
	cond1    interface{}
	operator string
	cond2    interface{}
	args     []interface{}
}

// add cond
func (w *Where) addCond(operator string, cond interface{}, args ...interface{}) *Where {
	w.operator = operator
	w.cond2 = cond
	w.args = append(w.args, args...)

	w2 := &Where{cond1: w}
	return w2
}

// And command
func (w *Where) And(cond interface{}, args ...interface{}) *Where {
	return w.addCond("AND", cond, args...)
}

// Or command
func (w *Where) Or(cond interface{}, args ...interface{}) *Where {
	return w.addCond("OR", cond, args...)
}

// ToString Where command
// ex) a=1 AND b=2
//
// return params
// - query string
// - args  []interface{}
func (w *Where) ToString() (string, []interface{}) {
	var args []interface{}
	query := ""

	if w1, ok := w.cond1.(*Where); ok {
		query1, args1 := w1.ToString()

		query += query1
		args = append(args, args1...)
	}

	if query != "" && w.operator != "" {
		query += " " + w.operator + " "
	}

	if w1, ok := w.cond2.(string); ok {
		query += w1
		args = append(args, w.args...)
	}

	return query, args
}

// ToCommand Where command
// ex) WHERE a=1 AND b=2
func (w *Where) ToCommand() (string, []interface{}) {
	var args []interface{}

	str, wargs := w.ToString()

	querystr := "WHERE " + str
	args = append(args, wargs...)

	return querystr, args
}
