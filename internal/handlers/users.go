package handlers

import (
	"cofee-shop-mongo/internal/auth"
	"cofee-shop-mongo/internal/utils"
	"cofee-shop-mongo/models"
	"context"
	"errors"
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
	mux.HandleFunc("POST /users", auth.WithJWTAuth(models.AdminAccess, h.createUser))
	mux.HandleFunc("POST /users/", auth.WithJWTAuth(models.AdminAccess, h.createUser))

	mux.HandleFunc("GET /users", auth.WithJWTAuth(models.AdminAccess, h.getAllUsers))
	mux.HandleFunc("GET /users/", auth.WithJWTAuth(models.AdminAccess, h.getAllUsers))

	mux.HandleFunc("GET /users/{id}", auth.WithJWTAuth(models.AdminAccess, h.getUserById))
	mux.HandleFunc("GET /users/{id}/", auth.WithJWTAuth(models.AdminAccess, h.getUserById))

	mux.HandleFunc("PUT /users/{id}", auth.WithJWTAuth(models.AdminAccess, h.updateUserById))
	mux.HandleFunc("PUT /users/{id}/", auth.WithJWTAuth(models.AdminAccess, h.updateUserById))
	
	mux.HandleFunc("DELETE /users/{id}", auth.WithJWTAuth(models.AdminAccess, h.deleteUserById))
	mux.HandleFunc("DELETE /users/{id}/", auth.WithJWTAuth(models.AdminAccess, h.deleteUserById))
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid request payload"))
		return
	}

	if err := validateUser(user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	id, err := h.Service.CreateUser(r.Context(), user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("could not create user, please try again later"))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully", "id": id})
}

func (h *UserHandler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("could not retrieve users, please try again later"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) getUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	user, err := h.Service.GetUserById(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("user not found"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) updateUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updatedUser models.User

	if err := utils.ParseJSON(r, &updatedUser); err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid request payload"))
		return
	}

	if err := validateUser(updatedUser); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.Service.UpdateUserById(r.Context(), id, updatedUser); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("could not update user, please try again later"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

func (h *UserHandler) deleteUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.Service.DeleteUserById(r.Context(), id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("could not delete user, please try again later"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func validateUser(user models.User) error {
	if user.UserID == "" {
		return errors.New("user ID cannot be empty")
	}
	if user.Username == "" {
		return errors.New("username cannot be empty")
	}
	if user.Email == "" {
		return errors.New("email cannot be empty")
	}
	if user.Password == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}
