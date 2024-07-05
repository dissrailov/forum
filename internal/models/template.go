package models

type TemplateData struct {
	Post            *Post
	Posts           *[]Post
	Comments        *[]Comment
	CurrentYear     int
	Form            any
	IsAuthenticated bool
	User            *User
}
