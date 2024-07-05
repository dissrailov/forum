package service

import (
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
		return nil, err
	}
	return post, nil
}

func (s *service) DislikePost(userID, postID int) error {
	return s.repo.DislikePost(userID, postID)
}

func (s *service) LikePost(userID, postID int) error {
	return s.repo.LikePost(userID, postID)
}

func (s *service) AddComment(postId, userId int, content string) error {
	err := s.repo.AddComment(postId, userId, content)
	return err
}

func (s *service) GetCommentByPostId(postId int) ([]models.Comment, error) {
	comment, err := s.repo.GetCommentByPostId(postId)
	if err != nil {
		return nil, err
	}
	return comment, nil
}
