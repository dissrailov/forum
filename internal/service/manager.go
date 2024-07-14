package service

import (
	"forum/internal/models"
	"forum/internal/repo"
	"net/http"
)

type service struct {
	repo repo.RepoI
}

type ServiceI interface {
	PostServiceI
	UserServiceI
}

func NewService(repo repo.RepoI) ServiceI {
	return &service{repo}
}

type PostServiceI interface {
	CreatePost(cookie string, form models.PostCreateForm, data *models.TemplateData) (*models.TemplateData, int, error)
	GetPostId(id int) (*models.Post, error)
	GetLastPost() ([]models.Post, error)
	DislikePost(userID, postID int) error
	LikePost(userID, postID int) error
	RemoveReaction(userID, postID int) error
	GetUserReaction(userID, postID int) (int, error)
	AddComment(postId, userId int, content string) error
	GetCommentByPostId(postId int) ([]models.Comment, error)
}

type UserServiceI interface {
	CreateUser(form models.UserSignupForm, data *models.TemplateData) (*models.TemplateData, error)
	Authenticate(form *models.UserLoginForm, data *models.TemplateData) (*models.Session, *models.TemplateData, error)
	DeleteSession(token string) error
	GetUser(r *http.Request) (*models.User, error)
	GetPassword(userId int) (string, error)
	ValidatePasswordForm(form *models.AccountPasswordUpdateForm) error
	UpdatePassword(userID int, oldPassword, newPassword string) error
}
