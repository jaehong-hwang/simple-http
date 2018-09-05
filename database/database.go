package database

import (
	"context"
	"database/sql"
)

// Pool struct
type Pool struct {
	Env          *Env
	MaxOpenConns int
	MaxIdleConns int
	SQLDB        *sql.DB
	context      context.Context
}

// NewPool func
func NewPool(env *Env, maxOpenConns, maxIdleConns int) (*Pool, error) {
	var err error
	db := &Pool{
		Env:          env,
		MaxOpenConns: maxOpenConns,
		MaxIdleConns: maxIdleConns,
		context:      context.Background(),
	}

	db.SQLDB, err = sql.Open(db.Env.Drive, db.Env.GetDataSourceName())

	if err != nil {
		return nil, err
	}

	if err = db.SQLDB.Ping(); err != nil {
		return nil, err
	}

	db.SQLDB.SetMaxOpenConns(db.MaxOpenConns)
	db.SQLDB.SetMaxIdleConns(db.MaxIdleConns)

	return db, nil
}

// Open DB Pool
func (p *Pool) Open() (*Connection, error) {
	var err error

	conn, err := p.SQLDB.Conn(p.context)
	if err != nil {
		return nil, err
	}

	return &Connection{
		conn:    conn,
		context: &p.context,
	}, nil
}

// Close pool
func (p *Pool) Close() error {
	return p.SQLDB.Close()
}
