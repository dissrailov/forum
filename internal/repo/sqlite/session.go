package sqlite

import (
	"fmt"
	"forum/internal/models"
)

func (s *Sqlite) CreateSession(session *models.Session) error {
	op := "sqlite.CreateSession"
	stmt := `INSERT INTO sessions(id, token, exp_time) VALUES(?, ?, ?)`
	_, err := s.DB.Exec(stmt, session.ID, session.Token, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
