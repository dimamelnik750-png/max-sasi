package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"maxsasi/internal/repository"
	"maxsasi/internal/service"
	"maxsasi/internal/user"
	"maxsasi/pkg/httpjson"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpjson.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var input user.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	u, err := h.authService.Register(input)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, u)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpjson.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var input user.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	tokens, err := h.authService.Login(input)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpjson.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	tokens, err := h.authService.Refresh(body.RefreshToken)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) handleAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		httpjson.WriteError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, service.ErrUsernameRequired),
		errors.Is(err, service.ErrPasswordRequired):
		httpjson.WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, repository.ErrUserAlreadyExists):
		httpjson.WriteError(w, http.StatusConflict, err.Error())
	default:
		httpjson.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
