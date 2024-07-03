package handlers

import (
	"errors"
	"forum/internal/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (h *HandlerApp) AccountView(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/account/view" {
		h.NotFound(w)
		return
	}
	h.methodResolver(w, r, h.AccountViewGet, h.AccountViewPost)
}

func (h *HandlerApp) AccountViewGet(w http.ResponseWriter, r *http.Request) {
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

func (h *HandlerApp) AccountViewPost(w http.ResponseWriter, r *http.Request) {
	oldPassword := r.FormValue("old_password")
	newpassword := r.FormValue("new_password")
	confirmpassword := r.FormValue("confirm_password")

	if newpassword != confirmpassword {
		http.Error(w, "New password do not match", http.StatusBadRequest)
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
	err = h.service.UpdatePassword(userid.ID, newpassword)
	if err != nil {
		http.Error(w, "Error updating password ", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}
