package database

// QueryError struct
type QueryError struct {
	QueryString string
	Parameters  []interface{}
	Message     string
}

func (e QueryError) Error() string {
	return e.Message
}
