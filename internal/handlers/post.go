package handlers

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"net/http"
	"strconv"
)

func (h *HandlerApp) postCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodPost)
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := h.service.CreatePost(title, content, expires)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", id), http.StatusSeeOther)
}

func (h *HandlerApp) postView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		h.NotFound(w)
		return
	}

	post, err := h.service.GetPostId(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.NotFound(w)
		} else {
			h.ServerError(w, err)
		}
		return
	}

	h.Render(w, http.StatusOK, "view.tmpl", &models.TemplateData{
		Post: post,
	})
}
