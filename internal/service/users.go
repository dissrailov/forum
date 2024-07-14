package service

import (
	"errors"
	"forum/internal/models"
	"forum/internal/pkg/cookie"
	"forum/internal/pkg/validator"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *service) CreateUser(form models.UserSignupForm, data *models.TemplateData) (*models.TemplateData, error) {
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	if !form.Valid() {
		data.Form = form
		return data, models.ErrNotValidPostForm
	}
	data.Form = form
	err := s.repo.CreateUser(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldErrors("email", "Email address is already in use")
			data.Form = form
			return data, models.ErrNotValidPostForm
		} else if errors.Is(err, models.ErrDuplicateName) {
			form.AddFieldErrors("name", "Usernmae is already in use")
			data.Form = form
			return data, models.ErrNotValidPostForm
		} else {
			return nil, err
		}
	}
	return data, err
}

func (s *service) Exists(id int) (bool, error) {
	return s.repo.Exists(id)
}

func (s *service) DeleteSession(token string) error {
	if err := s.repo.DeleteSessionByToken(token); err != nil {
		return err
	}
	return nil
}

func (s *service) Authenticate(form *models.UserLoginForm, data *models.TemplateData) (*models.Session, *models.TemplateData, error) {
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data.Form = form
		return nil, data, models.ErrNotValidPostForm
	}

	userId, err := s.repo.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldErrors("email", "Email or password is incorrect")
			data.Form = form
			return nil, data, models.ErrNotValidPostForm
		} else {
			return nil, nil, err
		}
	}
	session := models.NewSession(userId)
	if err = s.repo.DeleteSessionById(userId); err != nil {
		return nil, data, err
	}
	err = s.repo.CreateSession(session)
	if err != nil {
		return nil, data, err
	}
	return session, data, nil
}

func (s *service) GetUser(r *http.Request) (*models.User, error) {
	token := cookie.GetSessionCookie("session_id", r)
	userID, err := s.repo.GetUserIDByToken(token.Value)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(userID)
}

func (s *service) GetPassword(userId int) (string, error) {
	return s.repo.GetPassword(userId)
}

func (s *service) UpdatePassword(form models.AccountPasswordUpdateForm, data *models.TemplateData, userID int) (*models.TemplateData, error) {
	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "This field cannot be blank")
	form.CheckField(form.NewPassword == form.NewPasswordConfirmation, "newPasswordConfirmation", "Passwords do not match")
	if !form.Valid() {
		data.Form = form
		return data, nil
	}
	hashedPassword, err := s.GetPassword(userID)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(form.CurrentPassword))
	if err != nil {
		return nil, err
	}
	data.Form = form
	return data, err
}
func (s *service) GetUserReaction(userID, postID int) (int, error) {
	return s.repo.GetUserReaction(userID, postID)
}

func (s *service) RemoveReaction(userID, postID int) error {
	return s.repo.RemoveReaction(userID, postID)
}
