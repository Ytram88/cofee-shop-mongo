package handlers

import (
	"cofee-shop-mongo/internal/utils"
	"cofee-shop-mongo/models"
	"context"
	"net/http"
)

type AuthService interface {
	RegisterUser(ctx context.Context, payload models.RegisterUserPayload) error
	LoginUser(ctx context.Context, payload models.UserLoginPayload) (string, error)
}

type AuthHandler struct {
	Service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service}
}

func (h *AuthHandler) RegisterEndpoints(mux *http.ServeMux) {
	http.HandleFunc("POST /login", h.LoginUser)
	http.HandleFunc("POST /register", h.RegisterUser)
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// load user from request
	var userPayload models.UserLoginPayload
	if err := utils.ParseJSON(r, &userPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// get the token from service
	// hanlde errors better here
	token, err := h.Service.LoginUser(r.Context(), userPayload)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
	}

	//send token back
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})

}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userPayload models.RegisterUserPayload
	if err := utils.ParseJSON(r, &userPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.Service.RegisterUser(r.Context(), userPayload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
