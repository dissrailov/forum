package service

import (
	"forum/internal/models"
	"forum/internal/pkg/validator"
)

func (s *service) CreatePost(cookie string, form models.PostCreateForm, data *models.TemplateData) (*models.TemplateData, int, error) {
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Content, 250), "content", "This field cannot be more than 250 characters long")
	form.CheckField(validator.MaxChars(form.Title, 250), "title", "This field cannot be more than 250 characters long")

	if !form.Valid() {
		data.Form = form
		return data, 0, models.ErrInvalidCredentials
	}
	userID, err := s.repo.GetUserIDByToken(cookie)
	if err != nil {
		return nil, 0, err
	}
	postID, err := s.repo.CreatePost(form.Title, form.Content, userID)
	if err != nil {
		return nil, 0, err
	}
	data.Form = form
	return data, postID, nil

}

func (s *service) GetPostId(id int) (*models.Post, error) {
	return s.repo.GetPostId(id)
}

func (s *service) GetLastPost() ([]models.Post, error) {
	posts, err := s.repo.GetLastPost()
	if err != nil {
		return nil, err
	}
	return posts, nil
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
