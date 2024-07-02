package service

import (
	"forum/internal/models"
	"forum/internal/pkg/cookie"
	"net/http"
)

func (s *service) CreateUser(name, email, password string) error {
	err := s.repo.CreateUser(name, email, password)
	return err
}

func (s *service) Exists(id int) (bool, error) {
	return s.repo.Exists(id)
}

func (s *service) DeleteSession(token string) error {
	if err := s.repo.DeleteSessionByToken(token); err != nil {
		return err
	}
	return nil
}

func (s *service) Authenticate(email, password string) (*models.Session, error) {
	userId, err := s.repo.Authenticate(email, password)
	if err != nil {
		return nil, err
	}
	session := models.NewSession(userId)

	if err = s.repo.DeleteSessionById(userId); err != nil {
		return nil, err
	}
	err = s.repo.CreateSession(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *service) GetUser(r *http.Request) (*models.User, error) {
	token := cookie.GetSessionCookie("session_id", r)
	userID, err := s.repo.GetUserIDByToken(token.Value)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(userID)
}
