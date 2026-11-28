package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTodos_Unauthorized(t *testing.T) {
	server := setupTest(t)
	rr := doRequest(server, http.MethodGet, "/todos", "", nil)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetTodos_Empty(t *testing.T) {
	server := setupTest(t)
	token := registerAndLogin(server)
	rr := doRequest(server, http.MethodGet, "/todos", token, nil)

	assert.Equal(t, http.StatusOK, rr.Code)

	var items []map[string]any
	json.NewDecoder(rr.Body).Decode(&items)
	assert.Empty(t, items)
}

func TestCreateTodo(t *testing.T) {
	server := setupTest(t)
	token := registerAndLogin(server)
	rr := doRequest(server, http.MethodPost, "/todos", token, map[string]string{
		"title": "Buy milk",
	})

	assert.Equal(t, http.StatusCreated, rr.Code)

	var item map[string]any
	json.NewDecoder(rr.Body).Decode(&item)
	assert.Equal(t, "Buy milk", item["title"])
	assert.Equal(t, false, item["completed"])
	assert.NotEmpty(t, item["id"])
}

func TestCreateTodo_EmptyTitle(t *testing.T) {
	server := setupTest(t)
	token := registerAndLogin(server)
	rr := doRequest(server, http.MethodPost, "/todos", token, map[string]string{
		"title": "",
	})

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetTodoByID(t *testing.T) {
	server := setupTest(t)
	token := registerAndLogin(server)

	// створюємо todo
	rr := doRequest(server, http.MethodPost, "/todos", token, map[string]string{
		"title": "Buy milk",
	})
	var created map[string]any
	json.NewDecoder(rr.Body).Decode(&created)

	// отримуємо по id
	rr = doRequest(server, http.MethodGet, "/todos/"+created["id"].(string), token, nil)
	assert.Equal(t, http.StatusOK, rr.Code)

	var item map[string]any
	json.NewDecoder(rr.Body).Decode(&item)
	assert.Equal(t, created["id"], item["id"])
	assert.Equal(t, "Buy milk", item["title"])
}

func TestGetTodoByID_NotFound(t *testing.T) {
	server := setupTest(t)
	token := registerAndLogin(server)
	rr := doRequest(server, http.MethodGet, "/todos/00000000-0000-0000-0000-000000000000", token, nil)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUpdateTodo(t *testing.T) {
	server := setupTest(t)
	token := registerAndLogin(server)

	// створюємо todo
	rr := doRequest(server, http.MethodPost, "/todos", token, map[string]string{
		"title": "Buy milk",
	})
	var created map[string]any
	json.NewDecoder(rr.Body).Decode(&created)

	// оновлюємо
	rr = doRequest(server, http.MethodPut, "/todos/"+created["id"].(string), token, map[string]any{
		"title":     "Buy milk and eggs",
		"completed": true,
	})
	assert.Equal(t, http.StatusOK, rr.Code)

	var updated map[string]any
	json.NewDecoder(rr.Body).Decode(&updated)
	assert.Equal(t, "Buy milk and eggs", updated["title"])
	assert.Equal(t, true, updated["completed"])
}

func TestDeleteTodo(t *testing.T) {
	server := setupTest(t)
	token := registerAndLogin(server)

	// створюємо todo
	rr := doRequest(server, http.MethodPost, "/todos", token, map[string]string{
		"title": "Buy milk",
	})
	var created map[string]any
	json.NewDecoder(rr.Body).Decode(&created)

	// видаляємо
	rr = doRequest(server, http.MethodDelete, "/todos/"+created["id"].(string), token, nil)
	assert.Equal(t, http.StatusNoContent, rr.Code)

	// перевіряємо що не існує
	rr = doRequest(server, http.MethodGet, "/todos/"+created["id"].(string), token, nil)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
