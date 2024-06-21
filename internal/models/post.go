package models

import (
	"forum/internal/pkg/validator"
	"time"
)

type Post struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type PostCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}
