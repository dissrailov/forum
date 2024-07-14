package handlers

import (
	"errors"
	"forum/internal/models"
	"net/http"
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
	form := models.AccountPasswordUpdateForm{
		CurrentPassword:         r.FormValue("currentPassword"),
		NewPassword:             r.FormValue("newPassword"),
		NewPasswordConfirmation: r.FormValue("newPasswordConfirmation"),
	}
	data := h.NewTemplateData(r)
	err := r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}
	userid, err := h.service.GetUser(r)
	if err != nil {
		h.ServerError(w, err)
	}
	data, err = h.service.UpdatePassword(form, data, userid.ID)
	if err != nil {
		http.Error(w, "Error updating password ", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}
