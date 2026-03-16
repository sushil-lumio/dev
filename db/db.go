package db

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	instance *sql.DB
	once     sync.Once
)

const schema = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL DEFAULT ''
	)
`

func GetDB(dbPath string) (*sql.DB, error) {
	var err error
	once.Do(func() {
		instance, err = open(dbPath)
	})
	return instance, err
}

func NewTestDB() (*sql.DB, error) {
	return open(":memory:")
}

func open(dbPath string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err := conn.Exec(schema); err != nil {
		return nil, err
	}
	return conn, nil
}
