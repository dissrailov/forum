package handlers

import (
	"forum/internal/models"
	"forum/internal/pkg/cookie"
	"net/http"
	"strings"
)

func (h *HandlerApp) UserSignup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/signup" {
		h.NotFound(w)
		return
	}
	h.methodResolver(w, r, h.UserSignupGet, h.UserSignupPost)
}

func (h *HandlerApp) UserSignupPost(w http.ResponseWriter, r *http.Request) {
	form := models.UserSignupForm{
		Name:     r.FormValue("name"),
		Email:    strings.ToLower(r.FormValue("email")),
		Password: r.FormValue("password"),
	}
	data, err := h.NewTemplateData(r)
	if err != nil {
		h.ServerError(w, err)
	}
	err = r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	data, err = h.service.CreateUser(form, data)
	if err != nil {
		if err == models.ErrNotValidPostForm {
			h.Render(w, http.StatusBadRequest, "signup.tmpl", data)
			return
		} else {
			h.ServerError(w, err)
			return
		}
	}
	http.Redirect(w, r, "/user/login", http.StatusSeeOther) // redirect to login page
}

func (h *HandlerApp) UserSignupGet(w http.ResponseWriter, r *http.Request) {
	data, err := h.NewTemplateData(r)
	if err != nil {
		h.ServerError(w, err)
	}
	data.Form = models.UserSignupForm{}
	h.Render(w, http.StatusOK, "signup.tmpl", data)
}

func (h *HandlerApp) userLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/login" {
		h.NotFound(w)
		return
	}
	h.methodResolver(w, r, h.userLoginGet, h.userLoginPost)
}

func (h *HandlerApp) userLoginPost(w http.ResponseWriter, r *http.Request) {
	form := models.UserLoginForm{
		Email:    strings.ToLower(r.FormValue("email")),
		Password: r.FormValue("password"),
	}
	data, err := h.NewTemplateData(r)
	if err != nil {
		h.ServerError(w, err)
	}
	err = r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}
	session, templateData, err := h.service.Authenticate(&form, data)
	if err != nil {
		if err == models.ErrNotValidPostForm {
			h.Render(w, http.StatusBadRequest, "login.tmpl", templateData)
			return
		}
		h.ServerError(w, err)
		return
	}
	cookie.SetSessionCookie("session_id", w, session.Token, session.ExpTime)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *HandlerApp) userLoginGet(w http.ResponseWriter, r *http.Request) {
	data, err := h.NewTemplateData(r)
	if err != nil {
		h.ServerError(w, err)
	}
	data.Form = models.UserLoginForm{}
	h.Render(w, http.StatusOK, "login.tmpl", data)
}

func (h *HandlerApp) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	c := cookie.GetSessionCookie("session_id", r)
	if c != nil {
		h.service.DeleteSession(c.Value)
		cookie.ExpireSessionCookie("session_id", w)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
