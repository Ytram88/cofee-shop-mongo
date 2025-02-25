package service

import (
	"cofee-shop-mongo/models"
	"context"
	"fmt"
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
	const op = "service.GetAllInventoryItems"
	items, err := s.Repo.GetAllInventoryItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return items, nil
}

func (s *InventoryService) GetInventoryItemById(ctx context.Context, InventoryId string) (models.InventoryItem, error) {
	const op = "service.GetInventoryItemById"
	item, err := s.Repo.GetInventoryItemById(ctx, InventoryId)
	if err != nil {
		return models.InventoryItem{}, fmt.Errorf("%s: %w", op, err)
	}
	return item, nil
}

func (s *InventoryService) DeleteInventoryItemById(ctx context.Context, InventoryId string) error {
	const op = "service.DeleteInventoryItemById"
	err := s.Repo.DeleteInventoryItemById(ctx, InventoryId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *InventoryService) UpdateInventoryItemById(ctx context.Context, InventoryId string, item models.InventoryItem) error {
	const op = "service.UpdateInventoryItemById"
	item.IngredientID = InventoryId
	err := s.Repo.UpdateInventoryItemById(ctx, InventoryId, item)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *InventoryService) CreateInventoryItem(ctx context.Context, item models.InventoryItem) (string, error) {
	const op = "service.CreateInventoryItem"
	id, err := s.Repo.CreateInventoryItem(ctx, item)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *InventoryService) DeductStock(ctx context.Context, ingredientID string, qty float64) error {
	const op = "service.DeductStock"

	item, err := s.Repo.GetInventoryItemById(ctx, ingredientID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if item.Quantity < qty {
		return fmt.Errorf("%s: %w", op, ErrNotEnoughStock)
	}

	item.Quantity -= qty
	if err := s.Repo.UpdateInventoryItemById(ctx, ingredientID, item); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *InventoryService) HasSufficientStock(ctx context.Context, ingredientID string, requiredQty float64) bool {
	item, err := s.Repo.GetInventoryItemById(ctx, ingredientID)
	if err != nil {
		return false
	}
	return item.Quantity >= requiredQty
}
