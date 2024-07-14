package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/pkg/cookie"
	"log"
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

	form := models.PostCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
	}

	cookies := cookie.GetSessionCookie("session_id", r)
	data := h.NewTemplateData(r)

	data, id, err := h.service.CreatePost(cookies.Value, form, data)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	if len(form.FieldErrors) > 0 {
		data = h.NewTemplateData(r)
		data.Form = form
		h.Render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", id), http.StatusSeeOther)
}

func (h *HandlerApp) postCreateGet(w http.ResponseWriter, r *http.Request) {
	data := h.NewTemplateData(r)
	data.Form = models.PostCreateForm{}

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

	data := h.NewTemplateData(r)
	data.Post = post
	data.Comments = &comments
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
		log.Println(err)
		h.NotFound(w)
		return
	}

	userID, _ := h.service.GetUser(r) // Функция для получения идентификатора текущего пользователя

	reaction, err := h.service.GetUserReaction(userID.ID, postID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		h.ServerError(w, err)
		return
	}

	if reaction == 1 {
		if err := h.service.RemoveReaction(userID.ID, postID); err != nil {
			h.ServerError(w, err)
			return
		}
	} else {
		if reaction == -1 {
			if err := h.service.RemoveReaction(userID.ID, postID); err != nil {
				h.ServerError(w, err)
				return
			}
		}
		if err := h.service.LikePost(userID.ID, postID); err != nil {
			http.Error(w, "Failed to like post", http.StatusInternalServerError)
			return
		}
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
		log.Println(err)
		h.NotFound(w)
		return
	}

	userID, _ := h.service.GetUser(r) // Функция для получения идентификатора текущего пользователя

	reaction, err := h.service.GetUserReaction(userID.ID, postID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		h.ServerError(w, err)
		return
	}

	if reaction == -1 {
		if err := h.service.RemoveReaction(userID.ID, postID); err != nil {
			h.ServerError(w, err)
			return
		}
	} else {
		if reaction == 1 {
			if err := h.service.RemoveReaction(userID.ID, postID); err != nil {
				h.ServerError(w, err)
				return
			}
		}
		if err := h.service.DislikePost(userID.ID, postID); err != nil {
			http.Error(w, "Failed to dislike post", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *HandlerApp) AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postIDStr := r.FormValue("PostId")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil || postID < 1 {
			log.Println(err)
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
			fmt.Println(err)
			http.Error(w, "Unable to add comment", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", postID), http.StatusSeeOther)
	}
}
