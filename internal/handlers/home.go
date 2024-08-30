package handlers

import (
	"forum/internal/models"
	"net/http"
	"strconv"
)

func (h *HandlerApp) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.NotFound(w)
		return
	}

	categoryIDStr := r.URL.Query().Get("category")

	var categoryID int
	var err error

	if categoryIDStr != "" {
		categoryID, err = strconv.Atoi(categoryIDStr)
		if err != nil {
			h.ServerError(w, err)
			return
		}
	}

	var posts []models.Post

	if categoryID > 0 {
		posts, err = h.service.GetPostByCategory(categoryID)
	} else {
		posts, err = h.service.GetAllPosts()
	}
	if err != nil {
		h.ServerError(w, err)
		return
	}

	categories, err := h.service.GetAllCategories()
	if err != nil {
		h.ServerError(w, err)
		return
	}

	data, err := h.NewTemplateData(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	data.Posts = &posts
	data.Categories = &categories

	h.Render(w, http.StatusOK, "home.tmpl", data)
}
