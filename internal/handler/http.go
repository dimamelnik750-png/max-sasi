package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"maxsasi/internal/repository"
	"maxsasi/internal/service"
	"maxsasi/internal/todo"
	"maxsasi/pkg/httpjson"
)

type Handler struct {
	todoService service.TodoService
}

func New(todoService service.TodoService) *Handler {
	return &Handler{todoService: todoService}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		httpjson.WriteError(w, http.StatusNotFound, "route not found")
		return
	}

	if r.Method != http.MethodGet {
		httpjson.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	w.Write([]byte("todo API is running"))
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpjson.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *Handler) Todos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := h.todoService.GetAllTodos()
		if err != nil {
			h.handleServiceError(w, err)
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, items)
	case http.MethodPost:
		var input todo.CreateTodoInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
			return
		}

		item, err := h.todoService.CreateTodo(input)
		if err != nil {
			h.handleServiceError(w, err)
			return
		}

		httpjson.WriteJSON(w, http.StatusCreated, item)
	default:
		httpjson.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) TodoByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/todos/")
	if id == "" || id == r.URL.Path {
		httpjson.WriteError(w, http.StatusNotFound, "todo not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		item, err := h.todoService.GetTodoByID(id)
		if err != nil {
			h.handleServiceError(w, err)
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, item)
	case http.MethodPut:
		var input todo.UpdateTodoInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
			return
		}

		item, err := h.todoService.UpdateTodo(id, input)
		if err != nil {
			h.handleServiceError(w, err)
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, item)
	case http.MethodDelete:
		if err := h.todoService.DeleteTodo(id); err != nil {
			h.handleServiceError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	default:
		httpjson.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrTodoNotFound):
		httpjson.WriteError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrTitleRequired):
		httpjson.WriteError(w, http.StatusBadRequest, err.Error())
	default:
		httpjson.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
