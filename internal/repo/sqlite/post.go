package sqlite

import (
	"database/sql"
	"errors"
	"forum/internal/models"
)

func (s *Sqlite) CreatePost(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO posts (title, content, created, expires)
	VALUES (?, ?, datetime('now'), datetime('now', '+' || ? || ' day'))`
	result, err := s.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *Sqlite) GetPostId(id int) (*models.Post, error) {
	stmt := `SELECT id, title, content, created, expires FROM posts
WHERE expires > DATETIME('now') AND id = ?`
	row := s.DB.QueryRow(stmt, id)
	p := &models.Post{}
	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return p, nil
}

func (s *Sqlite) GetLastPost() (*[]models.Post, error) {
	stmt := `SELECT id, title, content, created, expires FROM posts
	WHERE expires > datetime('now') ORDER BY id DESC LIMIT 10`
	rows, err := s.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Expires)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &posts, nil
}
