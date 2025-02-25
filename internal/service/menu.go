package service

import (
	"cofee-shop-mongo/models"
	"context"
	"fmt"
)

type MenuRepository interface {
	CreateMenuItem(ctx context.Context, item models.MenuItem) (string, error)
	GetAllMenuItems(ctx context.Context) ([]models.MenuItem, error)
	GetMenuItemById(ctx context.Context, MenuId string) (models.MenuItem, error)
	DeleteMenuItemById(ctx context.Context, id string) error
	UpdateMenuItemById(ctx context.Context, id string, item models.MenuItem) error
}

type MenuService struct {
	Repo MenuRepository
}

func NewMenuService(repo MenuRepository) *MenuService {
	return &MenuService{Repo: repo}
}

func (s *MenuService) CreateMenuItem(ctx context.Context, item models.MenuItem) (string, error) {
	const op = "service.CreateMenuItem"
	id, err := s.Repo.CreateMenuItem(ctx, item)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *MenuService) GetAllMenuItems(ctx context.Context) ([]models.MenuItem, error) {
	const op = "service.GetAllMenuItems"

	items, err := s.Repo.GetAllMenuItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return items, nil
}

func (s *MenuService) GetMenuItemById(ctx context.Context, id string) (models.MenuItem, error) {
	const op = "service.GetMenuItemById"

	item, err := s.Repo.GetMenuItemById(ctx, id)
	if err != nil {
		return models.MenuItem{}, fmt.Errorf("%s: %w", op, err)
	}

	return item, nil
}

func (s *MenuService) DeleteMenuItemById(ctx context.Context, id string) error {
	const op = "service.DeleteMenuItemById"

	err := s.Repo.DeleteMenuItemById(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *MenuService) UpdateMenuItemById(ctx context.Context, id string, item models.MenuItem) error {
	const op = "service.UpdateMenuItemById"
	item.ProductId = id
	err := s.Repo.UpdateMenuItemById(ctx, id, item)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
