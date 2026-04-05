package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Handler struct {
	todoService TodoService
}

func NewHandler(todoService TodoService) *Handler {
	return &Handler{todoService: todoService}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		writeError(w, http.StatusNotFound, "route not found")
		return
	}

	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	w.Write([]byte("todo API is running"))
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *Handler) Todos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		todos, err := h.todoService.GetAllTodos()
		if err != nil {
			h.handleServiceError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, todos)
	case http.MethodPost:
		var input CreateTodoInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}

		todo, err := h.todoService.CreateTodo(input)
		if err != nil {
			h.handleServiceError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, todo)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) TodoByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/todos/")
	if id == "" || id == r.URL.Path {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		todo, err := h.todoService.GetTodoByID(id)
		if err != nil {
			h.handleServiceError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, todo)
	case http.MethodPut:
		var input UpdateTodoInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}

		todo, err := h.todoService.UpdateTodo(id, input)
		if err != nil {
			h.handleServiceError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, todo)
	case http.MethodDelete:
		if err := h.todoService.DeleteTodo(id); err != nil {
			h.handleServiceError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTodoNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrTitleRequired):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}
