package handlers

import (
	"errors"
	"forum/internal/models"
	"forum/internal/pkg/validator"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (h *HandlerApp) AccountView(w http.ResponseWriter, r *http.Request) {
	userId, err := h.service.GetUser(r)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		} else {
			h.ServerError(w, err)
		}
		return
	}
	data := h.NewTemplateData(r)
	data.User = userId
	h.Render(w, http.StatusOK, "account.tmpl", data)
}

func (h *HandlerApp) AccountChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/account/password" {
		h.NotFound(w)
		return
	}
	h.methodResolver(w, r, h.AccountChangePasswordGet, h.AccountChangePasswordPost)
}

func (h *HandlerApp) AccountChangePasswordGet(w http.ResponseWriter, r *http.Request) {
	data := h.NewTemplateData(r)
	data.Form = models.AccountPasswordUpdateForm{}
	h.Render(w, http.StatusOK, "password.tmpl", data)
}

func (h *HandlerApp) AccountChangePasswordPost(w http.ResponseWriter, r *http.Request) {
	oldPassword := r.FormValue("currentPassword")
	newPassword := r.FormValue("newPassword")
	newpasswordconfirm := r.FormValue("newPasswordConfirmation")
	form := models.AccountPasswordUpdateForm{
		CurrentPassword:         oldPassword,
		NewPassword:             newPassword,
		NewPasswordConfirmation: newpasswordconfirm,
	}
	err := r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "This field cannot be blank")
	form.CheckField(form.NewPassword == form.NewPasswordConfirmation, "newPasswordConfirmation", "Passwords do not match")

	if !form.Valid() {
		data := h.NewTemplateData(r)
		data.Form = form
		h.Render(w, http.StatusUnprocessableEntity, "password.tmpl", data)
		return
	}

	userid, err := h.service.GetUser(r)
	if err != nil {
		h.ServerError(w, err)
	}
	hashedPassword, err := h.service.GetPassword(userid.ID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(oldPassword))
	if err != nil {
		http.Error(w, "Old password incorrect", http.StatusUnauthorized)
		return
	}
	err = h.service.UpdatePassword(userid.ID, newPassword)
	if err != nil {
		http.Error(w, "Error updating password ", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}
