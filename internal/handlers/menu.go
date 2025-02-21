package handlers

import (
	"cofee-shop-mongo/internal/utils"
	"cofee-shop-mongo/models"
	"context"
	"log/slog"
	"net/http"
)

type MenuService interface {
	CreateMenuItem(ctx context.Context, item models.MenuItem) (string, error)
	GetAllMenuItems(ctx context.Context) ([]models.MenuItem, error)
	GetMenuItemById(ctx context.Context, id string) (models.MenuItem, error)
	UpdateMenuItemById(ctx context.Context, id string, item models.MenuItem) error
	DeleteMenuItemById(ctx context.Context, id string) error
}

type MenuHandler struct {
	Service MenuService
	Logger  *slog.Logger
}

func NewMenuHandler(service MenuService, logger *slog.Logger) *MenuHandler {
	return &MenuHandler{service, logger}
}

func (h *MenuHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /menu", h.createMenuItem)
	mux.HandleFunc("POST /menu/", h.createMenuItem)

	mux.HandleFunc("GET /menu", h.getAllMenu)
	mux.HandleFunc("GET /menu/", h.getAllMenu)

	mux.HandleFunc("GET /menu/{id}", h.getMenuItemById)
	mux.HandleFunc("GET /menu/{id}/", h.getMenuItemById)

	mux.HandleFunc("PUT /menu/{id}", h.updateMenuItemByIdById)
	mux.HandleFunc("PUT /menu/{id}/", h.updateMenuItemByIdById)

	mux.HandleFunc("DELETE /menu/{id}", h.deleteMenuItemById)
	mux.HandleFunc("DELETE /menu/{id}/", h.deleteMenuItemById)
}

func (h *MenuHandler) createMenuItem(w http.ResponseWriter, r *http.Request) {
	var item models.MenuItem

	err := utils.ParseJSON(r, &item)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err = h.Service.CreateMenuItem(r.Context(), item)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.Write([]byte("new item was created successfully"))
}
func (h *MenuHandler) getAllMenu(w http.ResponseWriter, r *http.Request) {
	items, err := h.Service.GetAllMenuItems(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, items)
}
func (h *MenuHandler) getMenuItemById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.Service.GetMenuItemById(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, item)
}
func (h *MenuHandler) updateMenuItemByIdById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updatedItem models.MenuItem

	err := utils.ParseJSON(r, &updatedItem)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.Service.UpdateMenuItemById(r.Context(), id, updatedItem)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, updatedItem)
}
func (h *MenuHandler) deleteMenuItemById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.Service.DeleteMenuItemById(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
