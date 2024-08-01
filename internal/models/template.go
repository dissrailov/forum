package models

type TemplateData struct {
	Post            *Post
	Posts           *[]Post
	PostsCreated    *Post
	PostsLiked      []Post
	Categories      *[]Category
	Comments        *[]Comment
	CurrentYear     int
	Form            any
	IsAuthenticated bool
	User            *User
}
