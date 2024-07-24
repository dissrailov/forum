package service

import "forum/internal/models"

func (s *service) GetCategoryByPostID(postID int) ([]models.Category, error) {
	category, err := s.repo.GetCategoryByPostID(postID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *service) GetAllCategories() ([]models.Category, error) {
	categories, err := s.repo.GetAllCategories()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *service) GetPostByCategory(categoryID int) ([]models.Post, error) {
	posts, err := s.repo.GetPostByCategory(categoryID)
	if err != nil {
		return nil, err
	}
	for i := range posts {
		categories, err := s.repo.GetCategoryByPostID(posts[i].ID)
		if err != nil {
			return nil, err
		}
		posts[i].Categories = categories
	}
	return posts, nil
}
