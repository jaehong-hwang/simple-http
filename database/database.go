package database

import (
	"database/sql"
)

// Pool struct
type Pool struct {
	Env          Env
	MaxOpenConns int
	MaxIdleConns int
	SQLDB        *sql.DB
}

// NewPool func
func NewPool(env Env, maxOpenConns, maxIdleConns int) (Pool, error) {
	db := Pool{
		Env:          env,
		MaxOpenConns: maxOpenConns,
		MaxIdleConns: maxIdleConns,
	}

	if err := db.Open(); err != nil {
		return Pool{}, err
	}

	db.Close()

	return db, nil
}

// Open DB Pool
func (p *Pool) Open() error {
	var err error
	p.SQLDB, err = sql.Open(p.Env.Drive, p.Env.GetDataSourceName())

	if err != nil {
		return err
	}

	if err = p.SQLDB.Ping(); err != nil {
		return err
	}

	p.SQLDB.SetMaxOpenConns(p.MaxOpenConns)
	p.SQLDB.SetMaxIdleConns(p.MaxIdleConns)
	return nil
}

// Close pool
func (p *Pool) Close() error {
	return p.SQLDB.Close()
}

// Table func
// query start with table name
func (p *Pool) Table(table string) *Query {
	query := &Query{connection: p}
	return query.From(table)
}

// Query start
func (p *Pool) Query() *Query {
	return &Query{connection: p}
}
