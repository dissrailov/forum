package handlers

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/pkg/validator"
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
	err := r.ParseForm()

	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := h.NewTemplateData(r)
		data.Form = form
		h.Render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = h.service.CreateUser(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldErrors("email", "Email address is already in use")
			data := h.NewTemplateData(r)
			data.Form = form
			h.Render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
			return
		} else {
			h.ServerError(w, err)
		}
		return
	}
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (h *HandlerApp) UserSignupGet(w http.ResponseWriter, r *http.Request) {
	data := h.NewTemplateData(r)
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
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (h *HandlerApp) userLoginGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user...")
}

func (h *HandlerApp) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
