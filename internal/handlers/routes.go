package handlers

import (
	"net/http"
)

func (h *HandlerApp) Routes() http.Handler {
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/post/view", h.postView)
	mux.HandleFunc("/post/create", h.postCreate)

	// users
	mux.HandleFunc("/user/signup", h.UserSignup)
	mux.HandleFunc("/user/login", (h.userLogin))
	mux.HandleFunc("/user/logout", h.userLogoutPost)

	return h.recoverPanic(h.Logrequest(SecureHeaders(mux)))
}
