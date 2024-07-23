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
		return nil, nil, err
	}
	err = s.repo.CreateSession(session)
	if err != nil {
		return nil, nil, err
	}
	data.Form = form
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

func (s *service) UpdatePassword(userID int, oldPassword, newPassword string) error {
	hashedPassword, err := s.GetPassword(userID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(oldPassword))
	if err != nil {
		return errors.New("old password incorrect")
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.repo.UpdatePassword(userID, string(newHashedPassword))
	if err != nil {
		return err
	}

	return nil
}

