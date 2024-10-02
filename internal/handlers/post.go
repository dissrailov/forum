package handlers

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/pkg/cookie"
	"net/http"
	"strconv"
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

	categories, err := h.service.GetAllCategories()
	if err != nil {
		h.ServerError(w, err)
		return
	}

	categoryIDsStr := r.Form["categoryIDs[]"]
	var categoryIDs []int
	for _, idStr := range categoryIDsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		categoryIDs = append(categoryIDs, id)
	}

	form := models.PostCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		CategoryIDs: categoryIDs,
	}
	cookies := cookie.GetSessionCookie("session_id", r)
	data, err := h.NewTemplateData(r)
	if err != nil {
		h.ServerError(w, err)
	}

	data.Categories = &categories

	data, id, err := h.service.CreatePost(cookies.Value, form, data)
	if err != nil {
		if err == models.ErrNotValidPostForm {
			h.Render(w, http.StatusBadRequest, "create.tmpl", data)
			return
		} else {
			h.ServerError(w, err)
			return
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", id), http.StatusSeeOther)
}

func (h *HandlerApp) postCreateGet(w http.ResponseWriter, r *http.Request) {
	data, err := h.NewTemplateData(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	_, err = h.service.GetUser(r)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintln("/user/login"), http.StatusSeeOther)
		return
	}
	categories, err := h.service.GetAllCategories()
	if err != nil {
		h.ServerError(w, err)
		return
	}
	data.Form = models.PostCreateForm{}
	data.Categories = &categories
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
	comments, err := h.service.GetCommentByPostId(id)
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
	}
	data.Post = post
	data.Comments = &comments
	data.Categories = &categories
	h.Render(w, http.StatusOK, "view.tmpl", data)
}

func (h *HandlerApp) LikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.FormValue("postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID < 1 {
		h.NotFound(w)
		return
	}

	userID, err := h.service.GetUser(r)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}
		h.ServerError(w, err)
	}

	err = h.service.LikePost(userID.ID, postID)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *HandlerApp) DislikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.FormValue("postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID < 1 {
		h.NotFound(w)
		return
	}

	userID, err := h.service.GetUser(r)
	if err != nil {
		h.ServerError(w, err)
	}

	err = h.service.DislikePost(userID.ID, postID)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *HandlerApp) AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postIDStr := r.FormValue("PostId")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil || postID < 1 {
			h.NotFound(w)
			return
		}
		userID, err := h.service.GetUser(r)
		content := r.FormValue("Content")
		if err != nil {
			h.ServerError(w, err)
			return
		}
		err = h.service.AddComment(postID, userID.ID, content)
		if err != nil {
			http.Error(w, "Unable to add comment", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", postID), http.StatusSeeOther)
	}
}

func (h *HandlerApp) LikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	commentIDStr := r.FormValue("commentID")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil || commentID < 1 {
		h.NotFound(w)
		return
	}

	userID, err := h.service.GetUser(r)
	if err != nil {
		h.ServerError(w, err)
	}

	err = h.service.LikeComment(userID.ID, commentID)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}

func (h *HandlerApp) DislikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	commentIDStr := r.FormValue("commentID")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil || commentID < 1 {
		h.NotFound(w)
		return
	}

	userID, err := h.service.GetUser(r)
	if err != nil {
		h.ServerError(w, err)
	}

	err = h.service.DislikePost(userID.ID, commentID)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}
