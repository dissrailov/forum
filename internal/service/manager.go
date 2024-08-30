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
	GetAllPosts() ([]models.Post, error)
	DislikePost(userID, postID int) error
	LikePost(userID, postID int) error
	GetLikedPosts(userID int) ([]models.Post, error)
	GetUserPosts(userID int) ([]models.Post, error)
	AddComment(postId, userId int, content string) error
	GetCategoryByPostID(postID int) ([]models.Category, error)
	GetAllCategories() ([]models.Category, error)
	GetCommentByPostId(postId int) ([]models.Comment, error)
	GetPostByCategory(categoryID int) ([]models.Post, error)
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
