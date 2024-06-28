package handlers

import (
	"fmt"
	"net/http"
)

func (h *HandlerApp) UserSignup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/signup" {
		h.NotFound(w)
		return
	}
	h.methodResolver(w, r, h.UserSignupGet, h.UserSignupPost)
}

func (h *HandlerApp) UserSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}

func (h *HandlerApp) UserSignupGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for signing up a new user...")
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
