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
	mux.HandleFunc("POST /menu", auth.WithJWTAuth(models.StaffAccess, h.createMenuItem))
	mux.HandleFunc("POST /menu/", auth.WithJWTAuth(models.StaffAccess, h.createMenuItem))

	mux.HandleFunc("GET /menu", auth.WithJWTAuth(models.StaffAccess, h.getAllMenuItems))
	mux.HandleFunc("GET /menu/", auth.WithJWTAuth(models.StaffAccess, h.getAllMenuItems))

	mux.HandleFunc("GET /menu/{id}", auth.WithJWTAuth(models.StaffAccess, h.getMenuItemById))
	mux.HandleFunc("GET /menu/{id}/", auth.WithJWTAuth(models.StaffAccess, h.getMenuItemById))

	mux.HandleFunc("PUT /menu/{id}", auth.WithJWTAuth(models.StaffAccess, h.updateMenuItemById))
	mux.HandleFunc("PUT /menu/{id}/", auth.WithJWTAuth(models.StaffAccess, h.updateMenuItemById))

	mux.HandleFunc("DELETE /menu/{id}", auth.WithJWTAuth(models.StaffAccess, h.deleteMenuItemById))
	mux.HandleFunc("DELETE /menu/{id}/", auth.WithJWTAuth(models.StaffAccess, h.deleteMenuItemById))
}

func (h *MenuHandler) createMenuItem(w http.ResponseWriter, r *http.Request) {
	var item models.MenuItem

	err := utils.ParseJSON(r, &item)
	if err != nil {
		h.Logger.Error("Failed to parse menu item request", "error", err)
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid request payload"))
		return
	}
	err = validateMenuItem(item)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	id, err := h.Service.CreateMenuItem(r.Context(), item)
	if err != nil {
		h.Logger.Error("Failed to create menu item", "error", err)
		utils.WriteError(w, http.StatusInternalServerError, errors.New("could not create menu item, please try again later"))
		return
	}

	h.Logger.Info("New menu item created", "id", id)
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "New menu item created successfully", "id": id})
}

func (h *MenuHandler) getAllMenuItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.Service.GetAllMenuItems(r.Context())
	if err != nil {
		h.Logger.Error("Failed to fetch menu items", "error", err)
		utils.WriteError(w, http.StatusInternalServerError, errors.New("could not retrieve menu items, please try again later"))
		return
	}

	h.Logger.Info("Fetched all menu items", "count", len(items))
	utils.WriteJSON(w, http.StatusOK, items)
}

func (h *MenuHandler) getMenuItemById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.Service.GetMenuItemById(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.Logger.Error("Menu item not found", "id", id, "error", err)
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("menu item \"%s\" not found", id))
		} else {
			h.Logger.Error("Failed to fetch menu item", "id", id, "error", err)
			utils.WriteError(w, http.StatusInternalServerError, errors.New("could not retrieve menu item, please try again later"))
		}
		return
	}

	h.Logger.Info("Fetched menu item", "id", id)
	utils.WriteJSON(w, http.StatusOK, item)
}

func (h *MenuHandler) updateMenuItemById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updatedItem models.MenuItem

	err := utils.ParseJSON(r, &updatedItem)
	if err != nil {
		h.Logger.Error("Failed to parse menu item update request", "id", id, "error", err)
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid request payload"))
		return
	}

	if updatedItem.ProductId == "" {
		updatedItem.ProductId = id
	}
	if updatedItem.ProductId != id {
		utils.WriteError(w, http.StatusBadRequest, errors.New("you cant change ProductId"))
	}
	if err = validateMenuItem(updatedItem); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	err = h.Service.UpdateMenuItemById(r.Context(), id, updatedItem)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.Logger.Error("Menu item not found", "id", id, "error", err)
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("menu item \"%s\" not found", id))
		} else {
			h.Logger.Error("Failed to update menu item", "id", id, "error", err)
			utils.WriteError(w, http.StatusInternalServerError, errors.New("could not update menu item, please try again later"))
		}
		return
	}

	h.Logger.Info("Updated menu item", "id", id)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Menu item updated successfully"})
}

func (h *MenuHandler) deleteMenuItemById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.Service.DeleteMenuItemById(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.Logger.Error("Menu item not found", "id", id, "error", err)
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("menu item \"%s\" not found", id))
		} else {
			h.Logger.Error("Failed to delete menu item", "id", id, "error", err)
			utils.WriteError(w, http.StatusInternalServerError, errors.New("could not delete menu item, please try again later"))
		}
		return
	}

	h.Logger.Info("Deleted menu item", "id", id)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Menu item deleted successfully"})
}

func validateMenuItem(item models.MenuItem) error {
	if item.ProductId == "" {
		return errors.New("product ID cannot be empty")
	}
	if item.Name == "" {
		return errors.New("name cannot be empty")
	}
	if item.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if len(item.Ingredients) == 0 {
		return errors.New("menu item must have at least one ingredient")
	}
	for _, ingredient := range item.Ingredients {
		if ingredient.IngredientID == "" {
			return errors.New("ingredient ID cannot be empty")
		}
		if ingredient.Quantity <= 0 {
			return errors.New("ingredient quantity must be greater than zero")
		}
	}
	return nil
}
