package mocks

import (
	"forum/internal/models"
	"testing"
)

type MockRepo struct{}

func NewMockRepo(t *testing.T) *MockRepo {
	return &MockRepo{}
}

func (r *MockRepo) AddCategory(postID int, categoryIDs []int) error {
	return nil
}

func (r *MockRepo) AddComment(postID, userID int, content string) error {
	return nil
}

func (r *MockRepo) Authenticate(email, password string) (int, error) {
	if email == "check@gmail.com" && password == "123check" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (r *MockRepo) CreateCategory(names []string) error {
	return nil
}

func (r *MockRepo) CreatePost(title, content string, userID int) (int, error) {
	return userID, nil
}

func (r *MockRepo) CreateSession(s *models.Session) error {
	return nil
}

func (r *MockRepo) CreateUser(name, email, password string) error {
	if name == "sam" && email == "check@gmail.com" {
		return nil
	}

	if email == "check@gmail.com" {
		return models.ErrDuplicateEmail
	}
	return nil
}

func (r *MockRepo) DeleteSessionById(userid int) error {
	return nil
}

func (r *MockRepo) DeleteSessionByToken(token string) error {
	return nil
}

func (r *MockRepo) DislikePost(userID, postID int) error {
	return nil
}

func (r *MockRepo) Exists(id int) (bool, error) {
	return true, nil
}

func (r *MockRepo) GetAllCategories() ([]models.Category, error) {
	// categories := []models.Category{
	// 	{ID: 1, Name: "Category1"},
	// 	{ID: 2, Name: "Category2"},
	// }

	return []models.Category{}, nil
}

func (r *MockRepo) GetAllPosts() ([]models.Post, error) {
	return []models.Post{}, nil
}

func (r *MockRepo) GetCategoryByPostID(postID int) ([]models.Category, error) {
	return []models.Category{}, nil
}

func (r *MockRepo) GetCommentByPostId(id int) ([]models.Comment, error) {
	return []models.Comment{}, nil
}

func (r *MockRepo) GetLikedPosts(userId int) ([]models.Post, error) {
	return []models.Post{{ID: 1, Title: "test Title", Content: "test Content"}}, nil
}

func (r *MockRepo) GetPassword(id int) (string, error) {
	return "123check", nil
}

func (r *MockRepo) GetPostByCategory(id int) ([]models.Post, error) {
	return []models.Post{}, nil
}

func (r *MockRepo) GetPostId(id int) (*models.Post, error) {
	if id == 1 {
		return &models.Post{ID: 1, Title: "test Title", Content: "test Content"}, nil
	}
	return nil, models.ErrNoRecord
}

func (r *MockRepo) GetUserByID(id int) (*models.User, error) {
	return &models.User{}, nil
}

func (r *MockRepo) GetUserIDByToken(token string) (int, error) {
	return 1, nil
}

func (r *MockRepo) GetUserPosts(id int) ([]models.Post, error) {
	return []models.Post{}, nil
}

func (r *MockRepo) GetUserReaction(id, postid int) (int, error) {
	return 1, nil
}

func (r *MockRepo) LikePost(id, postid int) error {
	return nil
}

func (r *MockRepo) RemoveReaction(id, postid int) error {
	return nil
}

func (r *MockRepo) UpdatePassword(id int, newPassword string) error {
	return nil
}
