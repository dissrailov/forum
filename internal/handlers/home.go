package handlers

import (
	"forum/internal/models"
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
	h.Render(w, http.StatusOK, "home.tmpl", &models.TemplateData{
		Posts: posts,
	})
}
