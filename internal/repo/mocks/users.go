package mocks

import (
	"forum/internal/models"
	"time"
)

type MockRepo struct{}

func (m MockRepo) Get(id int) (*models.User, error) {
	if id == 1 {
		u := &models.User{
			ID:      1,
			Name:    "Damir",
			Email:   "pearsayden@gmail.com",
			Created: time.Now(),
		}
		return u, nil
	}
	return nil, models.ErrNoRecord
}
