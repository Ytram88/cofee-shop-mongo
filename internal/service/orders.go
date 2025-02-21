package service

import (
	"cofee-shop-mongo/models"
	"context"
	"errors"
	"time"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, item models.Order) (string, error)
	GetAllOrders(ctx context.Context) ([]models.Order, error)
	GetOrderById(ctx context.Context, OrderId string) (models.Order, error)
	UpdateOrderById(ctx context.Context, OrderId string, item models.Order) error
	DeleteOrderById(ctx context.Context, OrderId string) error
}

type OrderService struct {
	OrderRepo        OrderRepository
	MenuService      *MenuService
	InventoryService *InventoryService
}

func NewOrderService(OrderRepo OrderRepository, MenuService *MenuService, InventoryService *InventoryService) *OrderService {
	return &OrderService{OrderRepo, MenuService, InventoryService}
}

func (s *OrderService) CreateOrder(ctx context.Context, order models.Order) (string, error) {
	if order.CustomerName == "" {
		return "", errors.New("customer name is required")
	}
	if len(order.Items) == 0 {
		return "", errors.New("order must have at least one item")
	}

	requiredIngredients := make(map[string]float64)
	for _, item := range order.Items {
		if item.Quantity <= 0 {
			return "", errors.New("invalid item quantity")
		}

		menuItem, err := s.MenuService.GetMenuItemById(ctx, item.ProductID)
		if err != nil {
			return "", errors.New("menu item not found: " + item.ProductID)
		}

		for _, ingredient := range menuItem.Ingredients {
			requiredIngredients[ingredient.IngredientID] += ingredient.Quantity * float64(item.Quantity)
		}
	}

	for ingredientID, needed := range requiredIngredients {
		invItem, err := s.InventoryService.GetInventoryItemById(ctx, ingredientID)
		if err != nil {
			return "", errors.New("inventory item not found: " + ingredientID)
		}
		if invItem.Quantity < needed {
			return "", errors.New("insufficient stock for ingredient: " + ingredientID)
		}
	}

	for ingredientID, needed := range requiredIngredients {
		invItem, _ := s.InventoryService.GetInventoryItemById(ctx, ingredientID)
		invItem.Quantity -= needed
		if err := s.InventoryService.UpdateInventoryItemById(ctx, ingredientID, invItem); err != nil {
			return "", errors.New("failed to update inventory for ingredient: " + ingredientID)
		}
	}

	order.Status = "open"
	order.CreatedAt = time.Now().Format(time.RFC3339)

	orderID, err := s.OrderRepo.CreateOrder(ctx, order)
	if err != nil {
		return "", errors.New("failed to create order")
	}

	return orderID, nil

}
func (s *OrderService) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	return s.OrderRepo.GetAllOrders(ctx)
}
func (s *OrderService) GetOrderById(ctx context.Context, OrderId string) (models.Order, error) {
	return s.OrderRepo.GetOrderById(ctx, OrderId)
}
func (s *OrderService) UpdateOrderById(ctx context.Context, OrderId string, order models.Order) error {
	if order.CustomerName == "" {
		return errors.New("customer name is required")
	}
	if len(order.Items) == 0 {
		return errors.New("order must have at least one item")
	}
	return s.OrderRepo.UpdateOrderById(ctx, OrderId, order)
}
func (s *OrderService) DeleteOrderById(ctx context.Context, OrderId string) error {
	return s.OrderRepo.DeleteOrderById(ctx, OrderId)
}
func (s *OrderService) CloseOrderById(ctx context.Context, OrderId string) error {
	order, err := s.OrderRepo.GetOrderById(ctx, OrderId)
	if err != nil {
		return errors.New("order not found")
	}

	if order.Status == "closed" {
		return errors.New("order already closed")
	}

	for _, item := range order.Items {
		menuItem, err := s.MenuService.GetMenuItemById(ctx, item.ProductID)
		if err != nil {
			return errors.New("menu item not found")
		}

		for _, ingredient := range menuItem.Ingredients {
			if !s.InventoryService.HasSufficientStock(ctx, ingredient.IngredientID, ingredient.Quantity*float64(item.Quantity)) {
				return errors.New("insufficient stock for ingredient: " + ingredient.IngredientID)
			}
		}
	}

	for _, item := range order.Items {
		menuItem, _ := s.MenuService.GetMenuItemById(ctx, item.ProductID)
		for _, ingredient := range menuItem.Ingredients {
			err = s.InventoryService.DeductStock(ctx, ingredient.IngredientID, ingredient.Quantity*float64(item.Quantity))
			if err != nil {
				return err
			}
		}
	}

	order.Status = "closed"
	return s.OrderRepo.UpdateOrderById(ctx, OrderId, order)
}
