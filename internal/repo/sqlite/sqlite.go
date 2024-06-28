package sqlite

import (
	"database/sql"
	"fmt"
)

type Sqlite struct {
	DB *sql.DB
}

func NewDB(dsn string) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	// Проверим соединение с базой данных
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	queries := []string{
		`CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`,
		`CREATE TABLE IF NOT EXISTS sessions (
		data BLOB NOT NULL,
		token INTEGER PRIMARY KEY,
		expiry TIMESTAMP NOT NULL
	);`,
		`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
  		name TEXT NOT NULL,
    	email TEXT NOT NULL,
    	hashed_password CHAR(60) NOT NULL,
    	created DATETIME NOT NULL
	);`,
	}
	for _, query := range queries {
		stmt, err := db.Prepare(query)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", "storage.sqlite.New", err)
		}
		_, err = stmt.Exec()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", "storage.sqlite.New", err)
		}
		stmt.Close()
	}
	return &Sqlite{DB: db}, nil
}
