package service

import (
	"errors"
	"forum/internal/models"
	"forum/internal/pkg/validator"
)

func (s *service) ValidatePasswordForm(form *models.AccountPasswordUpdateForm) error {
	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "This field cannot be blank")
	form.CheckField(form.NewPassword == form.NewPasswordConfirmation, "newPasswordConfirmation", "Passwords do not match")

	if !form.Valid() {
		return errors.New("invalid form data")
	}
	return nil
}
