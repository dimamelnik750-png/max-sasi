package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

var (
	todos  = make(map[string]Todo)
	nextID int
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var list []Todo

		for _, todo := range todos {
			list = append(list, todo)
		}

		writeJSON(w, http.StatusOK, list)

	case http.MethodPost:
		var req CreateTodoRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}

		if req.Title == "" {
			writeError(w, http.StatusBadRequest, "title is required")
			return
		}

		nextID++
		id := strconv.Itoa(nextID)

		todo := Todo{
			ID:        id,
			Title:     req.Title,
			Completed: false,
			CreatedAt: time.Now(),
		}

		todos[id] = todo
		writeJSON(w, http.StatusCreated, todo)

	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func todoByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/todos/"):]

	todo, exists := todos[id]
	if !exists {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, todo)

	case http.MethodPut:
		var req Todo

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}

		if req.Title == "" {
			writeError(w, http.StatusBadRequest, "title is required")
			return
		}

		todo.Title = req.Title
		todo.Completed = req.Completed
		todos[id] = todo

		writeJSON(w, http.StatusOK, todo)

	case http.MethodDelete:
		delete(todos, id)
		w.WriteHeader(http.StatusNoContent)

	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/todos", todosHandler)
	mux.HandleFunc("/todos/", todoByIDHandler)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
