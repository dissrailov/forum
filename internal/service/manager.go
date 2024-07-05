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
	CreatePost(title string, content string, expires int) (int, error)
	GetPostId(id int) (*models.Post, error)
	GetLastPost() (*[]models.Post, error)
	DislikePost(userID, postID int) error
	LikePost(userID, postID int) error
	RemoveReaction(userID, postID int) error
	GetUserReaction(userID, postID int) (int, error)
	AddComment(postId, userId int, content string) error
	GetCommentByPostId(postId int) ([]models.Comment, error)
}

type UserServiceI interface {
	CreateUser(name, email, password string) error
	Authenticate(email, password string) (*models.Session, error)
	DeleteSession(token string) error
	GetUser(r *http.Request) (*models.User, error)
	GetPassword(userId int) (string, error)
	UpdatePassword(userID int, hashedPassword string) error
}
