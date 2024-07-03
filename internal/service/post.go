package service

import (
	"fmt"
	"forum/internal/models"
)

func (s *service) CreatePost(title string, content string, expires int) (int, error) {
	return s.repo.CreatePost(title, content, expires)
}

func (s *service) GetPostId(id int) (*models.Post, error) {
	return s.repo.GetPostId(id)
}

func (s *service) GetLastPost() (*[]models.Post, error) {
	post, err := s.repo.GetLastPost()
	if err != nil {
		fmt.Println(err)
	}
	return post, nil
}

func (s *service) DislikePost(userID, postID int) error {
	return s.repo.DislikePost(userID, postID)
}

func (s *service) LikePost(userID, postID int) error {
	return s.repo.LikePost(userID, postID)
}
