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
	queries := `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(queries)
	if err != nil {
		fmt.Print(err)
	}
	return &Sqlite{DB: db}, nil
}
