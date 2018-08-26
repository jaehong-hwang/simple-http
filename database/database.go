package database

import (
	"database/sql"
)

// Connect to db func
func Connect(env Env, maxOpenConns, maxIdleConns int) (*sql.DB, error) {
	db, err := sql.Open(env.Drive, env.GetDataSourceName())

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	return db, nil
}
