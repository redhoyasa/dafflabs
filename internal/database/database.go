package database

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

const (
	Driver = "postgres"
)

type Database interface {
	Close() error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Conn struct {
	Database
}

func NewConn(databaseUrl string) (db Database, err error) {
	conn, err := sqlx.Connect(Driver, databaseUrl)

	if err != nil {
		return
	}

	db = Conn{conn}
	return
}
