package handlers

import (
	"errors"
	"forum/internal/models"
	"net/http"
)

func (h *HandlerApp) AccountView(w http.ResponseWriter, r *http.Request) {
	userId, err := h.service.GetUser(r)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		} else {
			h.ServerError(w, err)
		}
		return
	}
	data := h.NewTemplateData(r)
	data.User = userId
	h.Render(w, http.StatusOK, "account.tmpl", data)
}
