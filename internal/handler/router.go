package handler

import "net/http"

func NewRouter(h *Handler, auth *AuthHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/todos", h.Todos)
	mux.HandleFunc("/todos/", h.TodoByID)
	mux.HandleFunc("/auth/register", auth.Register)
	mux.HandleFunc("/auth/login", auth.Login)
	mux.HandleFunc("/auth/refresh", auth.Refresh)
	return mux
}
