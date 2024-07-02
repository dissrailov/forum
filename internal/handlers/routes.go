package handlers

import (
	"net/http"
)

func (h *HandlerApp) Routes() http.Handler {
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", h.Home)
	//post
	mux.HandleFunc("/post/view", h.RequireAuth(h.postView))
	mux.HandleFunc("/post/create", h.RequireAuth(h.postCreate))

	// users
	mux.HandleFunc("/user/signup", h.UserSignup)
	mux.HandleFunc("/user/login", h.userLogin)
	mux.HandleFunc("/user/logout", h.userLogoutPost)
	//account
	mux.HandleFunc("/account/view", h.RequireAuth(h.AccountView))

	return h.recoverPanic(h.Logrequest(SecureHeaders(mux)))
}
