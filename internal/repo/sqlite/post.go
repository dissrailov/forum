package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
)

func (s *Sqlite) CreatePost(title string, content string, userID int) (int, error) {
	op := "sqlite.CreatePost"
	stmt := `INSERT INTO posts (title, content, created, user_id)
	VALUES (?, ?, datetime('now'), ?)`
	result, err := s.DB.Exec(stmt, title, content, userID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return int(postID), nil
}

func (s *Sqlite) GetPostId(id int) (*models.Post, error) {
	op := "sqlite.GetPostId"
	stmt := `SELECT p.id, p.user_id, p.title, p.content, p.created, p.likes, p.dislikes, u.name
    FROM posts p
    JOIN users u ON p.user_id = u.id 
    WHERE p.id = ?`
	p := &models.Post{}
	err := s.DB.QueryRow(stmt, id).Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.Created, &p.Likes, &p.Dislikes, &p.UserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s : %w", op, models.ErrNoRecord)
		} else {
			return nil, fmt.Errorf("%s : %w", op, err)
		}
	}
	return p, nil
}

func (s *Sqlite) GetUserPosts(userID int) ([]models.Post, error) {
	op := "sqlite.GetUserPosts"
	query := `
        SELECT id, user_id, title, content, likes, dislikes, created
        FROM posts
        WHERE user_id = ?`
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.Likes,
			&post.Dislikes,
			&post.Created,
		); err != nil {
			return nil, fmt.Errorf("%s : %w", op, err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *Sqlite) GetAllPosts() ([]models.Post, error) {
	op := "sqlite.GetAllPosts"

	stmt := `SELECT id, title, content, likes, dislikes, created FROM posts
    ORDER BY id DESC`
	rows, err := s.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Likes, &p.Dislikes, &p.Created)
		if err != nil {
			return nil, fmt.Errorf("%s : %w", op, err)
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	return posts, nil
}

func (s *Sqlite) AddComment(postId, userId int, content string) error {
	op := "sqlite.AddComment"
	stmt, err := s.DB.Prepare("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	_, err = stmt.Exec(postId, userId, content)
	return nil
}

func (s *Sqlite) GetCommentByPostId(postId int) ([]models.Comment, error) {
	op := "sqlite.GetCommentByPostId"

	query := `
        SELECT c.id, c.post_id, c.user_id, u.name, c.content, c.created_at, c.likes, c.dislikes
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?
    `
	rows, err := s.DB.Query(query, postId)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	defer rows.Close()
	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostId, &comment.UserId, &comment.Username, &comment.Content, &comment.Created, &comment.Likes, &comment.Dislikes)
		if err != nil {
			return nil, fmt.Errorf("%s : %w", op, err)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
