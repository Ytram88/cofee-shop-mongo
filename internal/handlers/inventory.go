package handlers

import (
	"cofee-shop-mongo/internal/utils"
	"cofee-shop-mongo/models"
	"context"
	"log/slog"
	"net/http"
)

type InventoryService interface {
	CreateInventoryItem(ctx context.Context, item models.InventoryItem) (string, error)
	GetAllInventoryItems(ctx context.Context) ([]models.InventoryItem, error)
	GetInventoryItemById(ctx context.Context, InventoryId string) (models.InventoryItem, error)
	DeleteInventoryItemById(ctx context.Context, InventoryId string) error
	UpdateInventoryItemById(ctx context.Context, InventoryId string, item models.InventoryItem) error
}

type InventoryHandler struct {
	Service InventoryService
	Logger  *slog.Logger
}

func NewInventoryHandler(service InventoryService, logger *slog.Logger) *InventoryHandler {
	return &InventoryHandler{service, logger}
}

func (h *InventoryHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /inventory", h.createInventoryItem)
	mux.HandleFunc("POST /inventory/", h.createInventoryItem)

	mux.HandleFunc("GET /inventory", h.getAllInventoryItems)
	mux.HandleFunc("GET /inventory/", h.getAllInventoryItems)

	mux.HandleFunc("GET /inventory/{id}", h.getInventoryItemById)
	mux.HandleFunc("GET /inventory/{id}/", h.getInventoryItemById)

	mux.HandleFunc("PUT /inventory/{id}", h.updateInventoryItemById)
	mux.HandleFunc("PUT /inventory/{id}/", h.updateInventoryItemById)

	mux.HandleFunc("DELETE /inventory/{id}", h.deleteInventoryItemById)
	mux.HandleFunc("DELETE /inventory/{id}/", h.deleteInventoryItemById)
}

func (h *InventoryHandler) createInventoryItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem

	err := utils.ParseJSON(r, &item)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err = h.Service.CreateInventoryItem(r.Context(), item)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.Write([]byte("new item was created successfully"))
}
func (h *InventoryHandler) getAllInventoryItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.Service.GetAllInventoryItems(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, items)

}
func (h *InventoryHandler) getInventoryItemById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.Service.GetInventoryItemById(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, item)
}
func (h *InventoryHandler) updateInventoryItemById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updatedItem models.InventoryItem

	err := utils.ParseJSON(r, &updatedItem)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.Service.UpdateInventoryItemById(r.Context(), id, updatedItem)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, updatedItem)
}

func (h *InventoryHandler) deleteInventoryItemById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.Service.DeleteInventoryItemById(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
