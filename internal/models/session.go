package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        int
	Token     string
	ExpiresAt time.Time
}

func NewSession() *Session {
	return &Session{
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(100 * time.Minute),
	}
}
