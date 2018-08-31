package scope

// Operation struct
type Operation struct {
	val1     string
	operator string
	val2     string
}

// ToString operation
func (o *Operation) ToString() string {
	return o.val1 + o.operator + o.val2
}
