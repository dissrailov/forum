package models

import "time"

type AIResponse struct {
	ID               int
	PostID           int
	Content          string
	SimilarPostsJSON string
	SimilarPosts     []SimilarPost
	CreatedAt        time.Time
}

type SimilarPost struct {
	ID    int
	Title string
}
