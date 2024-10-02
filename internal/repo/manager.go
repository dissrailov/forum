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
	Category
}

type PostRepo interface {
	CreatePost(title string, content string, userID int) (int, error)
	GetPostId(id int) (*models.Post, error)
	GetUserPosts(userID int) ([]models.Post, error)
	GetAllPosts() ([]models.Post, error)
	LikePost(userID, postID int) error
	DislikePost(userID, postID int) error
	GetLikedPosts(userID int) ([]models.Post, error)
	GetUserReaction(userID, postID int) (int, error)
	RemoveReaction(userID, postID int) error
	AddComment(postId, userId int, content string) error
	GetCommentByPostId(postId int) ([]models.Comment, error)
	LikeComment(userID, commentID int) error
	DislikeComment(userID, commentID int) error
	RemoveReactionComm(userID, commentID int) error
	GetUserReactionComm(userID, commentID int) (int, error)
}

type Category interface {
	AddCategory(postID int, category []int) error
	CreateCategory(names []string) error
	GetCategoryByPostID(postID int) ([]models.Category, error)
	GetAllCategories() ([]models.Category, error)
	GetPostByCategory(categoryID int) ([]models.Post, error)
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
