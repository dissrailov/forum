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
	liked, err := h.service.GetLikedPosts(userId.ID)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	posts, err := h.service.GetUserPosts(userId.ID)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	data, err := h.NewTemplateData(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	data.LikedPosts = &liked
	data.UserPosts = &posts
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
	data, err := h.NewTemplateData(r)
	if err != nil {
		h.ServerError(w, err)
	}
	data.Form = models.AccountPasswordUpdateForm{}
	h.Render(w, http.StatusOK, "password.tmpl", data)
}

func (h *HandlerApp) AccountChangePasswordPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	form := models.AccountPasswordUpdateForm{
		CurrentPassword:         r.FormValue("currentPassword"),
		NewPassword:             r.FormValue("newPassword"),
		NewPasswordConfirmation: r.FormValue("newPasswordConfirmation"),
	}

	err = h.service.ValidatePasswordForm(&form)
	if err != nil {
		data, err := h.NewTemplateData(r)
		if err != nil {
			h.ClientError(w, http.StatusBadRequest)
			return
		}
		data.Form = form
		h.Render(w, http.StatusUnprocessableEntity, "password.tmpl", data)
		return
	}

	userID, err := h.service.GetUser(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	err = h.service.UpdatePassword(userID.ID, form.CurrentPassword, form.NewPassword)
	if err != nil {
		if err.Error() == "old password incorrect" {
			form.AddFieldErrors("currentPassword", "Old password is incorrect")
			data, err := h.NewTemplateData(r)
			if err != nil {
				h.ServerError(w, err)
				return
			}
			data.Form = form
			h.Render(w, http.StatusUnprocessableEntity, "password.tmpl", data)
		} else {
			h.ServerError(w, err)
		}
		return
	}

	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}
