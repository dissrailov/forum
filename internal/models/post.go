package models

import (
	"forum/internal/pkg/validator"
	"time"
)

type Post struct {
	ID         int
	UserID     int
	UserName   string
	Title      string
	Categories []Category
	Content    string
	ImageURL   string
	Likes      int
	Dislikes   int
	Created    time.Time
	Expires    time.Time
}

type PostCreateForm struct {
	Title               string
	Content             string
	ImageURL            string
	CategoryIDs         []int
	Category            []string
	validator.Validator `form:"-"`
}

type CommentForm struct {
	PostID   int
	UserID   int
	Content  string
	ImageURL string
	Token    string
	validator.Validator
}

type Comment struct {
	ID       int
	PostId   int
	UserId   int
	Username string
	Content  string
	ImageURL string
	Created  time.Time
	Likes    int
	Dislikes int
}
