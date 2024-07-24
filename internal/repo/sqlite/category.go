package sqlite

import (
	"fmt"
	"forum/internal/models"
)

func (s *Sqlite) CreateCategory(names []string) error {
	stmt := `INSERT INTO category (name) VALUES (?);`

	prepStmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return fmt.Errorf("fail to preprare %w", err)
	}
	defer prepStmt.Close()
	for _, name := range names {
		_, err := prepStmt.Exec(name)
		if err != nil {
			return fmt.Errorf("failed to execute: %w", err)
		}
	}
	return nil
}

func (s *Sqlite) GetCategoryByPostID(postID int) ([]models.Category, error) {
	op := "sqlite.GetCategoryByPostID"
	var categories []models.Category

	query := `SELECT c.id, c.name
              FROM category c
              JOIN post_category pc on c.id = pc.category_id
              WHERE pc.post_id = ?`
	rows, err := s.DB.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var category models.Category
		err = rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return categories, nil
}

func (s *Sqlite) AddCategory(postID int, categoryIDs []int) error {
	op := "sqlite.AddCategory"
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	for _, categoryID := range categoryIDs {
		_, err := stmt.Exec(postID, categoryID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// GetAllCategories retrieves all categories from the database.
func (s *Sqlite) GetAllCategories() ([]models.Category, error) {
	op := "sqlite.GetAllCategories"
	var categories []models.Category

	rows, err := s.DB.Query("SELECT id, name FROM category")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	fmt.Println(categories)
	return categories, nil
}

func (s *Sqlite) GetPostByCategory(categoryID int) ([]models.Post, error) {
	op := "sqlite.GetPostByCategory"
	rows, err := s.DB.Query(`
	SELECT p.id, p.title, p.content, p.user_id, p.created
	FROM posts p
	INNER JOIN post_category pc on p.id = pc.post_id
	WHERE pc.category_id = ?`, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Created)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}
