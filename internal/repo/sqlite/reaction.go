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


func (s *Sqlite) GetLikedPostsByUserID(userID int) ([]models.Post, error) {
	// Получаем все посты из базы данных
	rows, err := s.DB.Query("SELECT id, title, content, created FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created); err != nil {
			return nil, err
		}

		// Проверяем, лайкнул ли пользователь этот пост
		reaction, err := s.GetUserReaction(userID, post.ID)
		if err != nil {
			return nil, err
		}

		// Если реакция пользователя равна 1 (например, лайк), добавляем пост в результат
		if reaction == 1 {
			posts = append(posts, post)
		}
	}
	return posts, nil
}
