package handlers

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (h *HandlerApp) postCreate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		h.NotFound(w)
		return
	}
	h.methodResolver(w, r, h.postCreateGet, h.postCreatePost)
}

func (h *HandlerApp) postCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	form := models.PostCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(title) == "" {
		form.FieldErrors["title"] = "this field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		form.FieldErrors["title"] = "this field cannot be more than 100 characters long"
	}
	if strings.TrimSpace(content) == "" {
		form.FieldErrors["content"] = "this field cannot be blank"
	}
	if expires != 1 && expires != 7 && expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}
	if len(form.FieldErrors) > 0 {
		data := h.NewTemplateData(r)
		data.Form = form
		h.Render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := h.service.CreatePost(title, content, expires)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", id), http.StatusSeeOther)
}

func (h *HandlerApp) postCreateGet(w http.ResponseWriter, r *http.Request) {
	data := h.NewTemplateData(r)
	data.Form = models.PostCreateForm{
		Expires: 365,
	}
	h.Render(w, http.StatusOK, "create.tmpl", data)
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

	data := h.NewTemplateData(r)
	data.Post = post

	h.Render(w, http.StatusOK, "view.tmpl", data)

}
