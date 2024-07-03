package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"strings"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func (s *Sqlite) CreateUser(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created)
	         VALUES(?, ?, ?, datetime('now'));`
	_, err = s.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique && strings.Contains(sqliteErr.Error(), "email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (s *Sqlite) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, hashed_password FROM users WHERE email = ?`

	err := s.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (s *Sqlite) GetUserByID(id int) (*models.User, error) {
	op := "sqlite.GetUserByID"
	var u models.User
	stmt := `SELECT id, name, email, created FROM users WHERE id=?`
	err := s.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &u, nil
}

func (m *Sqlite) Exists(id int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?);`
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *Sqlite) GetPassword(userId int) (string, error) {
	var hashedPassword string

	err := m.DB.QueryRow("SELECT hashed_password FROM users WHERE id = ?", userId).Scan(&hashedPassword)
	if err != nil {
		return "", fmt.Errorf("failed to get hashed password: %v", err)
	}
	return hashedPassword, nil
}

func (m *Sqlite) UpdatePassword(userID int, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}
	_, err = m.DB.Exec("UPDATE users SET hashed_password = ? WHERE id = ?", hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("error updating password: %v", err)
	}
	return nil
}
