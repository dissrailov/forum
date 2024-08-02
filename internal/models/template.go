package models

type TemplateData struct {
	Post            *Post
	Posts           *[]Post
	Categories      *[]Category
	Comments        *[]Comment
	CurrentYear     int
	Form            any
	IsAuthenticated bool
	User            *User
	LikedPosts      *[]Post
	UserPosts       *[]Post
}
