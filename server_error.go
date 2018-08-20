package http

import "strconv"

// ServerError struct
type ServerError struct {
	Status int
}

func (e ServerError) Error() string {
	return strconv.Itoa(e.Status)
}
