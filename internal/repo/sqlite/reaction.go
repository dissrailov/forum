package sqlite

func (s *Sqlite) GetUserReaction(userID, postID int) (int, error) {
	var reaction int
	err := s.DB.QueryRow(`SELECT reaction FROM user_post_reactions WHERE user_id = ? AND post_id = ?`, userID, postID).Scan(&reaction)
	if err != nil {
		return 0, err
	}
	return reaction, nil
}

func (s *Sqlite) LikePost(userID, postID int) error {
	_, err := s.DB.Exec(`INSERT INTO user_post_reactions (user_id, post_id, reaction) VALUES (?, ?, 1)`, userID, postID)
	if err != nil {
		return err
	}
	_, err = s.DB.Exec(`UPDATE posts SET likes = likes + 1 WHERE id = ?`, postID)
	return err
}

func (s *Sqlite) DislikePost(userID, postID int) error {
	_, err := s.DB.Exec(`INSERT INTO user_post_reactions (user_id, post_id, reaction) VALUES (?, ?, -1)`, userID, postID)
	if err != nil {
		return err
	}
	_, err = s.DB.Exec(`UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?`, postID)
	return err
}

func (s *Sqlite) RemoveReaction(userID, postID int) error {
	var reaction int
	err := s.DB.QueryRow(`SELECT reaction FROM user_post_reactions WHERE user_id = ? AND post_id = ?`, userID, postID).Scan(&reaction)
	if err != nil {
		return err
	}
	_, err = s.DB.Exec(`DELETE FROM user_post_reactions WHERE user_id = ? AND post_id = ?`, userID, postID)
	if err != nil {
		return err
	}
	if reaction == 1 {
		_, err = s.DB.Exec(`UPDATE posts SET likes = likes - 1 WHERE id = ?`, postID)
	} else {
		_, err = s.DB.Exec(`UPDATE posts SET dislikes = dislikes - 1 WHERE id = ?`, postID)
	}
	return err
}
