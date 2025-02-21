package service

import (
	"cofee-shop-mongo/models"
	"context"
	"errors"
)

type InventoryRepository interface {
	GetAllInventoryItems(ctx context.Context) ([]models.InventoryItem, error)
	GetInventoryItemById(ctx context.Context, id string) (models.InventoryItem, error)
	DeleteInventoryItemById(ctx context.Context, id string) error
	UpdateInventoryItemById(ctx context.Context, id string, item models.InventoryItem) error
	CreateInventoryItem(ctx context.Context, item models.InventoryItem) (string, error)
}

type InventoryService struct {
	Repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) *InventoryService {
	return &InventoryService{Repo: repo}
}

func (s *InventoryService) GetAllInventoryItems(ctx context.Context) ([]models.InventoryItem, error) {
	return s.Repo.GetAllInventoryItems(ctx)
}

func (s *InventoryService) GetInventoryItemById(ctx context.Context, InventoryId string) (models.InventoryItem, error) {
	return s.Repo.GetInventoryItemById(ctx, InventoryId)
}

func (s *InventoryService) DeleteInventoryItemById(ctx context.Context, InventoryId string) error {
	return s.Repo.DeleteInventoryItemById(ctx, InventoryId)
}

func (s *InventoryService) UpdateInventoryItemById(ctx context.Context, InventoryId string, item models.InventoryItem) error {
	item.IngredientID = InventoryId
	return s.Repo.UpdateInventoryItemById(ctx, InventoryId, item)
}

func (s *InventoryService) CreateInventoryItem(ctx context.Context, item models.InventoryItem) (string, error) {
	return s.Repo.CreateInventoryItem(ctx, item)
}

func (s *InventoryService) DeductStock(ctx context.Context, ingredientID string, qty float64) error {
	item, err := s.Repo.GetInventoryItemById(ctx, ingredientID)
	if err != nil {
		return errors.New("inventory item not found")
	}

	if item.Quantity < qty {
		return errors.New("not enough stock for " + ingredientID)
	}

	item.Quantity -= qty
	return s.Repo.UpdateInventoryItemById(ctx, ingredientID, item)
}

func (s *InventoryService) HasSufficientStock(ctx context.Context, ingredientID string, requiredQty float64) bool {
	item, err := s.Repo.GetInventoryItemById(ctx, ingredientID)
	if err != nil {
		return false
	}
	return item.Quantity >= requiredQty
}
