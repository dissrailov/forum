package repo

import (
	"forum/internal/models"
	"forum/internal/repo/sqlite"
)

func NewRepo(dsn string) (RepoI, error) {
	return sqlite.NewDB(dsn)
}

type RepoI interface {
	PostRepo
}

type PostRepo interface {
	CreatePost(title string, content string, expires int) (int, error)
	GetPostId(id int) (*models.Post, error)
	GetLastPost() (*[]models.Post, error)
}
