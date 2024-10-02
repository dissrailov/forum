package sqlite

import (
	"fmt"
	"forum/internal/models"
)

func (s *Sqlite) GetUserReaction(userID, postID int) (int, error) {
	op := "sqlite.GetUserReaction"
	var reaction int
	err := s.DB.QueryRow(`SELECT reaction FROM user_post_reactions WHERE user_id = ? AND post_id = ?`, userID, postID).Scan(&reaction)
	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}
	return reaction, nil
}

func (s *Sqlite) LikePost(userID, postID int) error {
	op := "sqlite.LikePost"
	_, err := s.DB.Exec(`INSERT INTO user_post_reactions (user_id, post_id, reaction) VALUES (?, ?, 1)`, userID, postID)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	_, err = s.DB.Exec(`UPDATE posts SET likes = likes + 1 WHERE id = ?`, postID)
	return err
}

func (s *Sqlite) DislikePost(userID, postID int) error {
	op := "sqlite.DislikePost"
	_, err := s.DB.Exec(`INSERT INTO user_post_reactions (user_id, post_id, reaction) VALUES (?, ?, -1)`, userID, postID)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	_, err = s.DB.Exec(`UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?`, postID)
	return err
}

func (s *Sqlite) RemoveReaction(userID, postID int) error {
	op := "sqlite.RemoveReaction"
	var reaction int
	err := s.DB.QueryRow(`SELECT reaction FROM user_post_reactions WHERE user_id = ? AND post_id = ?`, userID, postID).Scan(&reaction)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	_, err = s.DB.Exec(`DELETE FROM user_post_reactions WHERE user_id = ? AND post_id = ?`, userID, postID)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	if reaction == 1 {
		_, err = s.DB.Exec(`UPDATE posts SET likes = likes - 1 WHERE id = ?`, postID)
	} else {
		_, err = s.DB.Exec(`UPDATE posts SET dislikes = dislikes - 1 WHERE id = ?`, postID)
	}
	return nil
}

func (s *Sqlite) GetLikedPosts(userID int) ([]models.Post, error) {
	op := "sqlite.GetLikedPosts"
	query := `
        SELECT p.id, p.title, p.content, p.created
        FROM posts p
        JOIN user_post_reactions upr ON p.id = upr.post_id
        WHERE upr.user_id = ? AND upr.reaction = 1` // Предполагаем, что 1 - это лайк

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created); err != nil {
			return nil, fmt.Errorf("%s : %w", op, err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *Sqlite) GetUserReactionComm(userID, commentID int) (int, error) {
	op := "sqlite.GetUserReactionComm"
	var reaction int
	err := s.DB.QueryRow(`SELECT reaction FROM user_comment_reactions WHERE user_id = ? AND comment_id = ?`, userID, commentID).Scan(&reaction)
	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}
	return reaction, nil
}

func (s *Sqlite) LikeComment(userID, commentID int) error {
	op := "sqlite.LikeComment"
	_, err := s.DB.Exec(`INSERT INTO user_comment_reactions (user_id, comment_id, reaction) VALUES (?, ?, 1)`, userID, commentID)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	_, err = s.DB.Exec(`UPDATE comments SET likes = likes + 1 WHERE id = ?`, commentID)
	return err
}

func (s *Sqlite) DislikeComment(userID, commentID int) error {
	op := "sqlite.DislikeComment"
	_, err := s.DB.Exec(`INSERT INTO user_comment_reactions (user_id, comment_id, reaction) VALUES (?, ?, -1)`, userID, commentID)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	_, err = s.DB.Exec(`UPDATE comments SET dislikes = dislikes + 1 WHERE id = ?`, commentID)
	return err
}

func (s *Sqlite) RemoveReactionComm(userID, commentID int) error {
	op := "sqlite.RemoveReactionComm"
	var reaction int
	err := s.DB.QueryRow(`SELECT reaction FROM user_comment_reactions WHERE user_id = ? AND comment_id = ?`, userID, commentID).Scan(&reaction)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	_, err = s.DB.Exec(`DELETE FROM user_comment_reactions WHERE user_id = ? AND comment_id = ?`, userID, commentID)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	if reaction == 1 {
		_, err = s.DB.Exec(`UPDATE comments SET likes = likes - 1 WHERE id = ?`, commentID)
	} else {
		_, err = s.DB.Exec(`UPDATE comments SET dislikes = dislikes - 1 WHERE id = ?`, commentID)
	}
	return nil
}
