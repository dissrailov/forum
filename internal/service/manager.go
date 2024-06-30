package service

import (
	"forum/internal/models"
	"forum/internal/repo"
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
}

type UserServiceI interface {
	CreateUser(name, email, password string) error
	Authenticate(email, password string) (*models.Session, error)
	DeleteSession(token string) error
}
