package handlers

import (
	"net/http"
)

func (h *HandlerApp) Routes() http.Handler {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/posts", h.apiGetPosts)
	mux.HandleFunc("/api/posts/view", h.apiGetPost)
	mux.HandleFunc("/api/posts/create", h.RequireAuthAPI(h.apiCreatePost))
	mux.HandleFunc("/api/posts/like", h.RequireAuthAPI(h.apiLikePost))
	mux.HandleFunc("/api/posts/dislike", h.RequireAuthAPI(h.apiDislikePost))
	mux.HandleFunc("/api/posts/ai", h.apiGetAIResponse)
	mux.HandleFunc("/api/posts/comments", h.apiAddComment)
	mux.HandleFunc("/api/comments/like", h.RequireAuthAPI(h.apiLikeComment))
	mux.HandleFunc("/api/comments/dislike", h.RequireAuthAPI(h.apiDislikeComment))
	mux.HandleFunc("/api/categories", h.apiGetCategories)
	mux.HandleFunc("/api/auth/me", h.apiGetMe)
	mux.HandleFunc("/api/auth/login", h.apiLogin)
	mux.HandleFunc("/api/auth/signup", h.apiSignup)
	mux.HandleFunc("/api/auth/logout", h.apiLogout)
	mux.HandleFunc("/api/account", h.RequireAuthAPI(h.apiGetAccount))
	mux.HandleFunc("/api/account/password", h.RequireAuthAPI(h.apiChangePassword))

	// Static images
	imgServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", imgServer))

	// Uploaded images
	uploads := http.FileServer(http.Dir("./uploads"))
	mux.Handle("/uploads/", http.StripPrefix("/uploads", uploads))

	// SPA catch-all — serves frontend/dist
	mux.HandleFunc("/", h.spaHandler("./frontend/dist"))

	return h.recoverPanic(h.Logrequest(h.CORS(mux)))
}
