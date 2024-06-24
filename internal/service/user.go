package service

import "forum/internal/models"

func (s *service) Authenticate() (*models.Session, error) {
	session := models.NewSession()
	if err := s.repo.CreateSession(session); err != nil {
		return nil, err
	}
	return session, nil
}
