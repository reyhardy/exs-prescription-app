package db

import (
	"context"
	"database/sql"

	"github.com/georgysavva/scany/v2/sqlscan"
)

type DB struct {
	conn *sql.DB
}

func New(driverName, dsn string) (*DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	// defer db.Close()
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{
		conn: db,
	}, nil
}

func (db *DB) Exec(ctx context.Context, q string, args ...any) error {
	if _, err := db.conn.ExecContext(ctx, q, args...); err != nil {
		return err
	}
	return nil
}

func (db *DB) Query(ctx context.Context, q string, dest any, args ...any) error {
	return sqlscan.Select(ctx, db.conn, dest, q, args...)
}

func (db *DB) QueryRow(ctx context.Context, q string, dest any, args ...any) error {
	return sqlscan.Get(ctx, db.conn, dest, q, args...)
}
