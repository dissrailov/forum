package handlers

import (
	"net/http"
)

func (h *HandlerApp) Routes() *http.ServeMux {
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/post/view", h.postView)
	mux.HandleFunc("/post/create", h.postCreate)
	return mux
}