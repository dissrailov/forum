package service

import (
	"database/sql"
	"errors"
	"forum/internal/models"
	"forum/internal/pkg/validator"
)

func (s *service) CreatePost(cookie string, form models.PostCreateForm, data *models.TemplateData) (*models.TemplateData, int, error) {
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Content, 250), "content", "This field cannot be more than 250 characters long")

	if !form.Valid() {
		data.Form = form
		return data, 0, models.ErrNotValidPostForm
	}

	userID, err := s.repo.GetUserIDByToken(cookie)
	if err != nil {
		return nil, 0, err
	}

	postID, err := s.repo.CreatePost(form.Title, form.Content, userID)
	if err != nil {
		return nil, 0, err
	}
	if err = s.repo.AddCategory(postID, form.CategoryIDs); err != nil {
		return nil, 0, err
	}
	data.Form = form
	return data, postID, nil
}

func (s *service) GetPostId(id int) (*models.Post, error) {
	posts, err := s.repo.GetPostId(id)
	if err != nil {
		return nil, err
	}

	categories, err := s.repo.GetCategoryByPostID(id)
	if err != nil {
		return nil, err
	}
	posts.Categories = categories
	return posts, nil
}

func (s *service) GetLastPost() ([]models.Post, error) {
	posts, err := s.repo.GetLastPost()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *service) DislikePost(userID, postID int) error {
	reaction, err := s.repo.GetUserReaction(userID, postID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if reaction == 1 || reaction == -1 {
		if err := s.repo.RemoveReaction(userID, postID); err != nil {
			return err
		}
	}

	if reaction != -1 {
		if err := s.repo.DislikePost(userID, postID); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) LikePost(userID, postID int) error {
	reaction, err := s.repo.GetUserReaction(userID, postID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if reaction == 1 || reaction == -1 {
		if err := s.repo.RemoveReaction(userID, postID); err != nil {
			return err
		}
	}

	if reaction != 1 {
		if err := s.repo.LikePost(userID, postID); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) GetLikedPostsByUserID(userID int) ([]models.Post, error) {
	reaction, err := s.repo.GetLikedPostsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return reaction, nil
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

func (s *service) GetLikedPosts(userID int) ([]models.Post, error) {
	liked, err := s.repo.GetLikedPosts(userID)
	if err != nil {
		return nil, err
	}
	return liked, nil
}

func (s *service) GetUserPosts(userID int) ([]models.Post, error) {
	posts, err := s.repo.GetUserPosts(userID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
