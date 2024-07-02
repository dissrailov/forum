package handlers

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/pkg/cookie"
	"net/http"
	"time"
)

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (h *HandlerApp) Logrequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (h *HandlerApp) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("connection", "close")
				h.ServerError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (h *HandlerApp) methodResolver(w http.ResponseWriter, r *http.Request, get, post func(w http.ResponseWriter, r *http.Request)) {
	switch r.Method {
	case http.MethodGet:
		get(w, r)
	case http.MethodPost:
		post(w, r)
	default:
		w.Header().Set("Content-Type", "text/plain")
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *HandlerApp) NewTemplateData(r *http.Request) *models.TemplateData {
	return &models.TemplateData{
		CurrentYear:     time.Now().Year(),
		IsAuthenticated: h.IsAuthenticated(r),
	}
}

func (h *HandlerApp) IsAuthenticated(r *http.Request) bool {
	cookie := cookie.GetSessionCookie("session_id", r)
	return cookie != nil && cookie.Value != ""
}

func (h *HandlerApp) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.IsAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

