package service

import (
	"cofee-shop-mongo/models"
	"context"
	"errors"
)

type MenuRepository interface {
	CreateMenuItem(ctx context.Context, item models.MenuItem) (string, error)
	GetAllMenuItems(ctx context.Context) ([]models.MenuItem, error)
	GetMenuItemById(ctx context.Context, id string) (models.MenuItem, error)
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
	if item.ProductId == "" {
		return "", errors.New("ProductId is empty")
	}
	if item.Name == "" {
		return "", errors.New("Name is empty")
	}
	if item.Ingredients == nil {
		return "", errors.New("Ingredients is empty")
	}
	return s.Repo.CreateMenuItem(ctx, item)
}

func (s *MenuService) GetAllMenuItems(ctx context.Context) ([]models.MenuItem, error) {
	return s.Repo.GetAllMenuItems(ctx)
}

func (s *MenuService) GetMenuItemById(ctx context.Context, id string) (models.MenuItem, error) {
	return s.Repo.GetMenuItemById(ctx, id)
}

func (s *MenuService) DeleteMenuItemById(ctx context.Context, id string) error {
	return s.Repo.DeleteMenuItemById(ctx, id)
}

func (s *MenuService) UpdateMenuItemById(ctx context.Context, id string, item models.MenuItem) error {
	item.ProductId = id
	if item.Ingredients == nil {
		return errors.New("Ingredients is empty")
	}
	return s.Repo.UpdateMenuItemById(ctx, id, item)
}
