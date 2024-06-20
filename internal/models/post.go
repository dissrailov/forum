package models

import "time"

type Post struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type PostCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}
