package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/pkg/cookie"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// JSON helper
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// saveUploadedFile saves an uploaded image to ./uploads/ and returns the URL path.
// Returns empty string if no file was uploaded.
func saveUploadedFile(r *http.Request, fieldName string) (string, error) {
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	ct := header.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "image/") {
		return "", fmt.Errorf("file is not an image")
	}

	if header.Size > 5<<20 { // 5MB
		return "", fmt.Errorf("file too large (max 5MB)")
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := uuid.New().String() + ext

	if err := os.MkdirAll("./uploads", 0o755); err != nil {
		return "", err
	}

	dst, err := os.Create(filepath.Join("./uploads", filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return "/uploads/" + filename, nil
}

// GET /api/posts?category=ID
func (h *HandlerApp) apiGetPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	categoryIDStr := r.URL.Query().Get("category")
	var posts []models.Post
	var err error

	if categoryIDStr != "" {
		categoryID, convErr := strconv.Atoi(categoryIDStr)
		if convErr != nil {
			writeError(w, http.StatusBadRequest, "invalid category ID")
			return
		}
		posts, err = h.service.GetPostByCategory(categoryID)
	} else {
		posts, err = h.service.GetAllPosts()
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get posts")
		h.ErrorLog.Println(err)
		return
	}

	writeJSON(w, http.StatusOK, posts)
}

// GET /api/posts/{id}
func (h *HandlerApp) apiGetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid post ID")
		return
	}

	post, err := h.service.GetPostId(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			writeError(w, http.StatusNotFound, "post not found")
		} else {
			writeError(w, http.StatusInternalServerError, "failed to get post")
			h.ErrorLog.Println(err)
		}
		return
	}

	comments, err := h.service.GetCommentByPostId(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get comments")
		h.ErrorLog.Println(err)
		return
	}

	categories, err := h.service.GetAllCategories()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get categories")
		h.ErrorLog.Println(err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"post":       post,
		"comments":   comments,
		"categories": categories,
	})
}

// POST /api/posts (multipart/form-data)
func (h *HandlerApp) apiCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid form data")
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryIDStrs := r.Form["categoryIDs"]
	var categoryIDs []int
	for _, s := range categoryIDStrs {
		id, err := strconv.Atoi(s)
		if err == nil {
			categoryIDs = append(categoryIDs, id)
		}
	}

	imageURL, err := saveUploadedFile(r, "image")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	form := models.PostCreateForm{
		Title:       title,
		Content:     content,
		ImageURL:    imageURL,
		CategoryIDs: categoryIDs,
	}

	c := cookie.GetSessionCookie("session_id", r)
	if c == nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	data := &models.TemplateData{}
	data, id, err := h.service.CreatePost(c.Value, form, data)
	if err != nil {
		if err == models.ErrNotValidPostForm {
			formData, _ := data.Form.(models.PostCreateForm)
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error":       "validation failed",
				"fieldErrors": formData.FieldErrors,
			})
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create post")
		h.ErrorLog.Println(err)
		return
	}

	if h.aiService != nil {
		go h.aiService.GenerateAndStore(id, form.Title, form.Content, form.ImageURL)
	}

	writeJSON(w, http.StatusCreated, map[string]any{"id": id})
}

// POST /api/posts/like?id=X
func (h *HandlerApp) apiLikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || postID < 1 {
		writeError(w, http.StatusBadRequest, "invalid post ID")
		return
	}

	user, err := h.service.GetUser(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	if err := h.service.LikePost(user.ID, postID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to like post")
		h.ErrorLog.Println(err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /api/posts/dislike?id=X
func (h *HandlerApp) apiDislikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || postID < 1 {
		writeError(w, http.StatusBadRequest, "invalid post ID")
		return
	}

	user, err := h.service.GetUser(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	if err := h.service.DislikePost(user.ID, postID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to dislike post")
		h.ErrorLog.Println(err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// GET /api/posts/ai?id=X
func (h *HandlerApp) apiGetAIResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid post ID")
		return
	}

	aiResp, err := h.service.GetAIResponse(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get AI response")
		h.ErrorLog.Println(err)
		return
	}

	if aiResp == nil {
		writeJSON(w, http.StatusOK, nil)
		return
	}

	writeJSON(w, http.StatusOK, aiResp)
}

// POST /api/posts/comments?id=X (multipart/form-data)
func (h *HandlerApp) apiAddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || postID < 1 {
		writeError(w, http.StatusBadRequest, "invalid post ID")
		return
	}

	user, err := h.service.GetUser(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid form data")
		return
	}

	content := r.FormValue("content")

	imageURL, err := saveUploadedFile(r, "image")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	form := models.CommentForm{Content: content, ImageURL: imageURL}
	data := &models.TemplateData{}
	_, err = h.service.AddComment(data, form, postID, user.ID)
	if err != nil {
		if err == models.ErrNotValidPostForm {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error":       "validation failed",
				"fieldErrors": map[string]string{"Content": "Comment cannot be blank or more than 100 characters"},
			})
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to add comment")
		h.ErrorLog.Println(err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

// POST /api/comments/like?id=X&postID=Y
func (h *HandlerApp) apiLikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || commentID < 1 {
		writeError(w, http.StatusBadRequest, "invalid comment ID")
		return
	}

	user, err := h.service.GetUser(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	if err := h.service.LikeComment(user.ID, commentID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to like comment")
		h.ErrorLog.Println(err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /api/comments/dislike?id=X
func (h *HandlerApp) apiDislikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || commentID < 1 {
		writeError(w, http.StatusBadRequest, "invalid comment ID")
		return
	}

	user, err := h.service.GetUser(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	if err := h.service.DislikeComment(user.ID, commentID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to dislike comment")
		h.ErrorLog.Println(err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// GET /api/categories
func (h *HandlerApp) apiGetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	categories, err := h.service.GetAllCategories()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get categories")
		h.ErrorLog.Println(err)
		return
	}

	writeJSON(w, http.StatusOK, categories)
}

// GET /api/auth/me
func (h *HandlerApp) apiGetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if !h.IsAuthenticated(r) {
		writeJSON(w, http.StatusOK, nil)
		return
	}

	user, err := h.service.GetUser(r)
	if err != nil {
		writeJSON(w, http.StatusOK, nil)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":      user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"created": user.Created,
	})
}

// POST /api/auth/login
func (h *HandlerApp) apiLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	form := &models.UserLoginForm{
		Email:    strings.ToLower(req.Email),
		Password: req.Password,
	}
	data := &models.TemplateData{}

	session, _, err := h.service.Authenticate(form, data)
	if err != nil {
		if err == models.ErrNotValidPostForm {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error":       "validation failed",
				"fieldErrors": form.FieldErrors,
			})
			return
		}
		writeError(w, http.StatusInternalServerError, "authentication failed")
		h.ErrorLog.Println(err)
		return
	}

	cookie.SetSessionCookie("session_id", w, session.Token, session.ExpTime)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /api/auth/signup
func (h *HandlerApp) apiSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	form := models.UserSignupForm{
		Name:     req.Name,
		Email:    strings.ToLower(req.Email),
		Password: req.Password,
	}
	data := &models.TemplateData{}

	_, err := h.service.CreateUser(form, data)
	if err != nil {
		if err == models.ErrNotValidPostForm {
			formData, _ := data.Form.(models.UserSignupForm)
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error":       "validation failed",
				"fieldErrors": formData.FieldErrors,
			})
			return
		}
		writeError(w, http.StatusInternalServerError, "signup failed")
		h.ErrorLog.Println(err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

// POST /api/auth/logout
func (h *HandlerApp) apiLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	c := cookie.GetSessionCookie("session_id", r)
	if c != nil {
		h.service.DeleteSession(c.Value)
		cookie.ExpireSessionCookie("session_id", w)
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// GET /api/account
func (h *HandlerApp) apiGetAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user, err := h.service.GetUser(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	liked, err := h.service.GetLikedPosts(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get liked posts")
		h.ErrorLog.Println(err)
		return
	}

	posts, err := h.service.GetUserPosts(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get user posts")
		h.ErrorLog.Println(err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user":       map[string]any{"id": user.ID, "name": user.Name, "email": user.Email, "created": user.Created},
		"likedPosts": liked,
		"userPosts":  posts,
	})
}

// POST /api/account/password
func (h *HandlerApp) apiChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		CurrentPassword         string `json:"currentPassword"`
		NewPassword             string `json:"newPassword"`
		NewPasswordConfirmation string `json:"newPasswordConfirmation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	form := &models.AccountPasswordUpdateForm{
		CurrentPassword:         req.CurrentPassword,
		NewPassword:             req.NewPassword,
		NewPasswordConfirmation: req.NewPasswordConfirmation,
	}

	if err := h.service.ValidatePasswordForm(form); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error":       "validation failed",
			"fieldErrors": form.FieldErrors,
		})
		return
	}

	user, err := h.service.GetUser(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	if err := h.service.UpdatePassword(user.ID, form.CurrentPassword, form.NewPassword); err != nil {
		if err.Error() == "old password incorrect" {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error":       "validation failed",
				"fieldErrors": map[string]string{"currentPassword": "Old password is incorrect"},
			})
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to change password")
		h.ErrorLog.Println(err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// RequireAuthAPI middleware for API routes
func (h *HandlerApp) RequireAuthAPI(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !h.IsAuthenticated(r) {
			writeError(w, http.StatusUnauthorized, "not authenticated")
			return
		}
		next.ServeHTTP(w, r)
	}
}

// CORS middleware for dev mode
func (h *HandlerApp) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// SPA handler - serves index.html for all non-API, non-static routes
func (h *HandlerApp) spaHandler(distDir string) http.HandlerFunc {
	fileServer := http.FileServer(http.Dir(distDir))
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Skip API routes
		if strings.HasPrefix(path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Try to serve static file first
		if info, err := http.Dir(distDir).Open(path); err == nil {
			stat, _ := info.Stat()
			info.Close()
			if !stat.IsDir() {
				fileServer.ServeHTTP(w, r)
				return
			}
		}

		// Fall back to index.html for SPA routing
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", distDir))
	}
}

