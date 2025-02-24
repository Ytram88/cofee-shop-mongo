package handlers

import (
	"cofee-shop-mongo/internal/utils"
	"cofee-shop-mongo/models"
	"context"
	"log/slog"
	"net/http"
)

type UserService interface {
	CreateUser(ctx context.Context, user models.User) (string, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, userId string) (models.User, error)
	UpdateUserById(ctx context.Context, userId string, user models.User) error
	DeleteUserById(ctx context.Context, userId string) error
}

type UserHandler struct {
	Service UserService
	Logger  *slog.Logger
}

func NewUserHandler(service UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{service, logger}
}

func (h *UserHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /users", h.createUser)
	mux.HandleFunc("GET /users", h.getAllUsers)
	mux.HandleFunc("GET /users/{id}", h.getUserById)
	mux.HandleFunc("PUT /users/{id}", h.updateUserById)
	mux.HandleFunc("DELETE /users/{id}", h.deleteUserById)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := utils.ParseJSON(r, &user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err = h.Service.CreateUser(r.Context(), user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func (h *UserHandler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) getUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	user, err := h.Service.GetUserById(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) updateUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updatedUser models.User

	err := utils.ParseJSON(r, &updatedUser)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.Service.UpdateUserById(r.Context(), id, updatedUser)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

func (h *UserHandler) deleteUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.Service.DeleteUserById(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
