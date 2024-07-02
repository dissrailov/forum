package sqlite

import (
	"fmt"
	"forum/internal/models"
)

func (s *Sqlite) CreateSession(session *models.Session) error {
	op := "sqlite.CreateSession"
	stmt := `INSERT INTO sessions(user_id, token, exp_time) VALUES(?, ?, ?)`
	_, err := s.DB.Exec(stmt, session.UserID, session.Token, session.ExpTime)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Sqlite) DeleteSessionById(userid int) error {
	op := "sqlite.DeleteSessionById"
	stmt := `DELETE FROM sessions WHERE user_id = ?`
	if _, err := s.DB.Exec(stmt, userid); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Sqlite) DeleteSessionByToken(token string) error {
	op := "sqlite.DeleteSessionByToken"
	stmt := `DELETE FROM sessions WHERE token = ?`
	if _, err := s.DB.Exec(stmt, token); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Sqlite) GetUserIDByToken(token string) (int, error) {
	op := "sqlite.GetUserIDByToken"
	stmt := `SELECT user_id FROM sessions WHERE token = ?`
	var userID int

	err := s.DB.QueryRow(stmt, token).Scan(&userID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}
