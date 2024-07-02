package models

type TemplateData struct {
	Post            *Post
	Posts           *[]Post
	CurrentYear     int
	Form            any
	IsAuthenticated bool
}
