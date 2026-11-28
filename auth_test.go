package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	server := setupTest(t)
	rr := doRequest(server, http.MethodPost, "/auth/register", "", map[string]string{
		"username": "testuser",
		"password": "testpass123",
	})

	assert.Equal(t, http.StatusCreated, rr.Code)

	var user map[string]any
	json.NewDecoder(rr.Body).Decode(&user)
	assert.Equal(t, "testuser", user["username"])
	assert.NotEmpty(t, user["id"])
}

func TestRegister_DuplicateUsername(t *testing.T) {
	server := setupTest(t)
	doRequest(server, http.MethodPost, "/auth/register", "", map[string]string{
		"username": "testuser",
		"password": "testpass123",
	})

	rr := doRequest(server, http.MethodPost, "/auth/register", "", map[string]string{
		"username": "testuser",
		"password": "testpass123",
	})

	assert.Equal(t, http.StatusConflict, rr.Code)
}

func TestLogin(t *testing.T) {
	server := setupTest(t)
	doRequest(server, http.MethodPost, "/auth/register", "", map[string]string{
		"username": "testuser",
		"password": "testpass123",
	})

	rr := doRequest(server, http.MethodPost, "/auth/login", "", map[string]string{
		"username": "testuser",
		"password": "testpass123",
	})

	assert.Equal(t, http.StatusOK, rr.Code)

	var tokens map[string]string
	json.NewDecoder(rr.Body).Decode(&tokens)
	assert.NotEmpty(t, tokens["access_token"])
	assert.NotEmpty(t, tokens["refresh_token"])
}

func TestLogin_WrongPassword(t *testing.T) {
	server := setupTest(t)
	doRequest(server, http.MethodPost, "/auth/register", "", map[string]string{
		"username": "testuser",
		"password": "testpass123",
	})

	rr := doRequest(server, http.MethodPost, "/auth/login", "", map[string]string{
		"username": "testuser",
		"password": "wrongpassword",
	})

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestRefresh(t *testing.T) {
	server := setupTest(t)
	doRequest(server, http.MethodPost, "/auth/register", "", map[string]string{
		"username": "testuser",
		"password": "testpass123",
	})

	rr := doRequest(server, http.MethodPost, "/auth/login", "", map[string]string{
		"username": "testuser",
		"password": "testpass123",
	})

	var tokens map[string]string
	json.NewDecoder(rr.Body).Decode(&tokens)

	rr = doRequest(server, http.MethodPost, "/auth/refresh", "", map[string]string{
		"refresh_token": tokens["refresh_token"],
	})

	assert.Equal(t, http.StatusOK, rr.Code)

	var newTokens map[string]string
	json.NewDecoder(rr.Body).Decode(&newTokens)
	assert.NotEmpty(t, newTokens["access_token"])
	assert.NotEmpty(t, newTokens["refresh_token"])
}
