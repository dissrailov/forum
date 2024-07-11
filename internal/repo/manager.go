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
	UserRepo
	SessionRepo
}

type PostRepo interface {
	CreatePost(title string, content string, userid, expires int) (int, error)
	GetPostId(id int) (*models.Post, error)
	GetLastPost() (*[]models.Post, error)
	LikePost(userID, postID int) error
	DislikePost(userID, postID int) error
	GetUserReaction(userID, postID int) (int, error)
	RemoveReaction(userID, postID int) error
	AddComment(postId, userId int, content string) error
	GetCommentByPostId(postId int) ([]models.Comment, error)
}

type UserRepo interface {
	CreateUser(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	GetUserByID(id int) (*models.User, error)
	GetPassword(userId int) (string, error)
	UpdatePassword(userID int, hashedPassword string) error
}

type SessionRepo interface {
	CreateSession(session *models.Session) error
	DeleteSessionById(userid int) error
	DeleteSessionByToken(token string) error
	GetUserIDByToken(token string) (int, error)
}
