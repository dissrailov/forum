package models

import "time"

type Post struct {
	ID       int
	Title    string
	Content  string
	Likes    int
	Dislikes int
	Created  time.Time
	Expires  time.Time
}

type PostCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

type UserPostReaction struct {
	UserID   int
	PostID   int
	Reaction int // 1 для лайка, -1 для дизлайка
}

//добавить в форму лайки
