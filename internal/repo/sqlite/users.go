package sqlite

import (
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

func (m *Sqlite) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *Sqlite) Exists(id int) (bool, error) {
	return false, nil
}
