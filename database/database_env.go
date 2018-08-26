package database

import "fmt"

// Env struct
type Env struct {
	Drive    string
	Host     string
	Database string
	User     string
	Password string
	Charset  string
}

// GetDataSourceName from environment variable
func (e *Env) GetDataSourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&autocommit=true", e.User, e.Password, e.Host, e.Database, e.Charset)
}
