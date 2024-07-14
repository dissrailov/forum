package handlers

import (
	"net/http"
)

func (h *HandlerApp) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.NotFound(w)
		return
	}
	posts, err := h.service.GetLastPost()
	if err != nil {
		h.ServerError(w, err)
		return
	}
	data := h.NewTemplateData(r)
	data.Posts = &posts // Присваиваем значение posts типа []models.Post
	h.Render(w, http.StatusOK, "home.tmpl", data)
}
