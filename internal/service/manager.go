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
}

func NewService(repo repo.RepoI) ServiceI {
	return &service{repo}
}

type PostServiceI interface {
	CreatePost(title string, content string, expires int) (int, error)
	GetPostId(id int) (*models.Post, error)
	GetLastPost() (*[]models.Post, error)
}
