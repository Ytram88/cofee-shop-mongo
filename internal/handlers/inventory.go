package handlers

import (
	"cofee-shop-mongo/internal/auth"
	"cofee-shop-mongo/internal/repository"
	"cofee-shop-mongo/internal/utils"
	"cofee-shop-mongo/models"
	"context"
	"errors"
	"fmt"
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
	mux.HandleFunc("POST /inventory", auth.WithJWTAuth(models.StaffAccess, h.createInventoryItem))
	mux.HandleFunc("POST /inventory/", auth.WithJWTAuth(models.StaffAccess, h.createInventoryItem))

	mux.HandleFunc("GET /inventory", auth.WithJWTAuth(models.StaffAccess, h.getAllInventoryItems))
	mux.HandleFunc("GET /inventory/", auth.WithJWTAuth(models.StaffAccess, h.getAllInventoryItems))

	mux.HandleFunc("GET /inventory/{id}", auth.WithJWTAuth(models.StaffAccess, h.getInventoryItemById))
	mux.HandleFunc("GET /inventory/{id}/", auth.WithJWTAuth(models.StaffAccess, h.getInventoryItemById))

	mux.HandleFunc("PUT /inventory/{id}", auth.WithJWTAuth(models.StaffAccess, h.updateInventoryItemById))
	mux.HandleFunc("PUT /inventory/{id}/", auth.WithJWTAuth(models.StaffAccess, h.updateInventoryItemById))

	mux.HandleFunc("DELETE /inventory/{id}", auth.WithJWTAuth(models.StaffAccess, h.deleteInventoryItemById))
	mux.HandleFunc("DELETE /inventory/{id}/", auth.WithJWTAuth(models.StaffAccess, h.deleteInventoryItemById))
}

func (h *InventoryHandler) createInventoryItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem

	err := utils.ParseJSON(r, &item)
	if err != nil {
		h.Logger.Error("Failed to parse inventory item request", "error", err)
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid request payload"))
		return
	}

	if err := validateInventoryItem(item); err != nil {
		h.Logger.Warn("Inventory item validation failed", "error", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	id, err := h.Service.CreateInventoryItem(r.Context(), item)
	if err != nil {
		h.Logger.Error("Failed to create inventory item", "error", err)
		utils.WriteError(w, http.StatusInternalServerError, errors.New("could not create inventory item, please try again later"))
		return
	}

	h.Logger.Info("New inventory item created", "id", id)
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "New inventory item created successfully", "id": id})
}

func (h *InventoryHandler) getAllInventoryItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.Service.GetAllInventoryItems(r.Context())
	if err != nil {
		h.Logger.Error("Failed to fetch inventory items", "error", err)
		utils.WriteError(w, http.StatusInternalServerError, errors.New("could not retrieve inventory items, please try again later"))
		return
	}

	h.Logger.Info("Fetched all inventory items", "count", len(items))
	utils.WriteJSON(w, http.StatusOK, items)
}

func (h *InventoryHandler) getInventoryItemById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.Service.GetInventoryItemById(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.Logger.Error("Inventory item not found", "id", id, "error", err)
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("inventory item \"%s\" not found", id))
		} else {
			h.Logger.Error("Failed to fetch inventory item", "id", id, "error", err)
			utils.WriteError(w, http.StatusNotFound, errors.New("inventory item not found"))
			return
		}
	}

	h.Logger.Info("Fetched inventory item", "id", id)
	utils.WriteJSON(w, http.StatusOK, item)
}

func (h *InventoryHandler) updateInventoryItemById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updatedItem models.InventoryItem

	err := utils.ParseJSON(r, &updatedItem)
	if err != nil {
		h.Logger.Error("Failed to parse inventory item update request", "id", id, "error", err)
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid request payload"))
		return
	}
	//validation user input
	if updatedItem.IngredientID == "" {
		updatedItem.IngredientID = id
	}
	if updatedItem.IngredientID != id {
		utils.WriteError(w, http.StatusBadRequest, errors.New("you cant change inventoryID"))
	}
	if err := validateInventoryItem(updatedItem); err != nil {
		h.Logger.Warn("Inventory item validation failed", "error", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.Service.UpdateInventoryItemById(r.Context(), id, updatedItem)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.Logger.Error("Inventory item not found", "id", id, "error", err)
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("inventory item \"%s\" not found", id))
		} else {
			h.Logger.Error("Failed to update inventory item", "id", id, "error", err)
			utils.WriteError(w, http.StatusInternalServerError, errors.New("could not update inventory item, please try again later"))
			return
		}
	}

	h.Logger.Info("Updated inventory item", "id", id)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Inventory item updated successfully"})
}

func (h *InventoryHandler) deleteInventoryItemById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.Service.DeleteInventoryItemById(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.Logger.Error("Inventory item not found", "id", id, "error", err)
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("inventory item \"%s\" not found", id))
		} else {
			h.Logger.Error("Failed to delete inventory item", "id", id, "error", err)
			utils.WriteError(w, http.StatusInternalServerError, errors.New("could not delete inventory item, please try again later"))
			return
		}
	}

	h.Logger.Info("Deleted inventory item", "id", id)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Inventory item deleted successfully"})
}

func validateInventoryItem(item models.InventoryItem) error {
	if item.IngredientID == "" {
		return errors.New("ingredient ID cannot be empty")
	}
	if item.Name == "" {
		return errors.New("name cannot be empty")
	}
	if item.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	if item.Unit == "" {
		return errors.New("unit cannot be empty")
	}
	return nil
}
