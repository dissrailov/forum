package models

import "time"

type Post struct {
	ID       int
	UserID   int
	UserName string
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

type Comment struct {
	ID       int
	PostId   int
	UserId   int
	Username string
	Content  string
	Created  time.Time
}
