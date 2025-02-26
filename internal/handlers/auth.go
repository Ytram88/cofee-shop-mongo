package handlers

import (
	"cofee-shop-mongo/internal/service"
	"cofee-shop-mongo/internal/utils"
	"cofee-shop-mongo/models"
	"context"
	"errors"
	"log/slog"
	"net/http"
)

type AuthService interface {
	RegisterUser(ctx context.Context, payload models.RegisterUserPayload) (string, error)
	LoginUser(ctx context.Context, payload models.UserLoginPayload) (string, error)
}

type AuthHandler struct {
	Service AuthService
	logger  *slog.Logger
}

func NewAuthHandler(service AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{service, logger}
}

func (h *AuthHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /login", h.LoginUser)
	mux.HandleFunc("POST /register", h.RegisterUser)
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Parse form values
	if err := r.ParseForm(); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Load user from request
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validate the payload
	if email == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("empty email"))
		return
	}
	if password == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("empty password"))
		return
	}

	// Create userPayload for service call
	userPayload := models.UserLoginPayload{
		Email:    email,
		Password: password,
	}

	// Get the token from service
	token, err := h.Service.LoginUser(r.Context(), userPayload)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	// Send token back
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Parse form values
	if err := r.ParseForm(); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Load user from request
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validate the payload
	if email == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("empty email"))
		return
	}
	if username == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("empty username"))
		return
	}
	if password == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("empty password"))
		return
	}

	// Create userPayload for service call
	userPayload := models.RegisterUserPayload{
		Email:    email,
		Username: username,
		Password: password,
	}

	// Start registering user
	userId, err := h.Service.RegisterUser(r.Context(), userPayload)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyExists) {
			utils.WriteError(w, http.StatusConflict, errors.New("user with this email already exists"))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	h.logger.Info("successfully registered user", slog.String("user id", userId))
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "user created successfully"})
}
