package service

import (
	"database/sql"
	"errors"
	"forum/internal/models"
	"forum/internal/pkg/validator"
	"strings"
)

func (s *service) CreatePost(cookie string, form models.PostCreateForm, data *models.TemplateData) (*models.TemplateData, int, error) {
	form.CheckField(validator.NotBlank(form.Title), "Title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "Title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "Content", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Content, 250), "Content", "This field cannot be more than 250 characters long")

	if !form.Valid() {
		data.Form = form
		return data, 0, models.ErrNotValidPostForm
	}

	userID, err := s.repo.GetUserIDByToken(cookie)
	if err != nil {
		return nil, 0, err
	}

	title := strings.TrimSpace(form.Title)
	content := strings.TrimSpace(form.Content)

	postID, err := s.repo.CreatePost(title, content, userID)
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

func (s *service) GetAllPosts() ([]models.Post, error) {
	posts, err := s.repo.GetAllPosts()
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

func (s *service) AddComment(data *models.TemplateData, form models.CommentForm, postID int, userId int) (*models.TemplateData, error) {
	form.CheckField(validator.NotBlank(form.Content), "Content", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Content, 100), "Content", "This field cannot be more than 100 characters long")
	content := strings.TrimSpace(form.Content)

	if !form.Valid() {
		data.Form = form
		return data, models.ErrNotValidPostForm
	}

	data.Form = form
	err := s.repo.AddComment(postID, userId, content)
	return data, err
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

func (s *service) DislikeComment(userID, commentID int) error {
	reaction, err := s.repo.GetUserReactionComm(userID, commentID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if reaction == 1 || reaction == -1 {
		if err := s.repo.RemoveReactionComm(userID, commentID); err != nil {
			return err
		}
	}

	if reaction != -1 {
		if err := s.repo.DislikeComment(userID, commentID); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) LikeComment(userID, commentID int) error {
	reaction, err := s.repo.GetUserReactionComm(userID, commentID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if reaction == 1 || reaction == -1 {
		if err := s.repo.RemoveReactionComm(userID, commentID); err != nil {
			return err
		}
	}

	if reaction != 1 {
		if err := s.repo.LikeComment(userID, commentID); err != nil {
			return err
		}
	}
	return nil
}
