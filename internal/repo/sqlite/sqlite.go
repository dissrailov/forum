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

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		likes INTEGER DEFAULT 0,
    	dislikes INTEGER DEFAULT 0,
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			token TEXT NOT NULL,
			exp_time TIMESTAMP NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS user_post_reactions (
			user_id INTEGER,
			post_id INTEGER,
			reaction INTEGER, -- 1 для лайка, -1 для дизлайка
			PRIMARY KEY (user_id, post_id)
		);`,
		`CREATE TABLE IF NOT EXISTS user_comment_reactions (
			user_id INTEGER,
			comment_id INTEGER,
			reaction INTEGER, -- 1 для лайка, -1 для дизлайка
			PRIMARY KEY (user_id, comment_id)
		);`,
		`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
  		name TEXT NOT NULL,
    	email TEXT NOT NULL,
    	hashed_password CHAR(60) NOT NULL,
    	created DATETIME NOT NULL,
    	UNIQUE(email)
	);`,
		`CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
	likes INTEGER DEFAULT 0,
    dislikes INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
	);`,
		`CREATE TABLE IF NOT EXISTS category (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);`,
		`CREATE TABLE IF NOT EXISTS post_category (
		category_id INTEGER,
		post_id INTEGER, 
		PRIMARY KEY (category_id, post_id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (category_id) REFERENCES category(id)
	);`,
	}
	for _, query := range queries {
		stmt, err := db.Prepare(query)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", "sqlite.NewDB", err)
		}
		_, err = stmt.Exec()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", "sqlite.NewDB", err)
		}
		stmt.Close()
	}
	return &Sqlite{DB: db}, nil
}
