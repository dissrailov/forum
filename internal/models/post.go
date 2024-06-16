package models

import "time"

type Post struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// TODO: realization post
