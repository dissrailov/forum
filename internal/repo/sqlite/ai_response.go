package sqlite

import (
	"database/sql"
	"forum/internal/models"
)

func (s *Sqlite) CreateAIResponse(postID int, content string, similarPostsJSON string) error {
	stmt := `INSERT INTO ai_responses (post_id, content, similar_posts, created_at)
	VALUES (?, ?, ?, datetime('now'))`
	_, err := s.DB.Exec(stmt, postID, content, similarPostsJSON)
	return err
}

func (s *Sqlite) GetAIResponseByPostID(postID int) (*models.AIResponse, error) {
	stmt := `SELECT id, post_id, content, similar_posts, created_at FROM ai_responses WHERE post_id = ?`
	r := &models.AIResponse{}
	err := s.DB.QueryRow(stmt, postID).Scan(&r.ID, &r.PostID, &r.Content, &r.SimilarPostsJSON, &r.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r, nil
}
