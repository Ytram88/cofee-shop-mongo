package service

import (
	"cofee-shop-mongo/models"
	"context"
	"fmt"
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
	const op = "service.CreateOrder"

	order.Status = "open"
	order.CreatedAt = time.Now().Format(time.RFC3339)

	orderID, err := s.OrderRepo.CreateOrder(ctx, order)
	if err != nil {
		return "", fmt.Errorf("%s: failed to create order, %w", op, err)
	}

	return orderID, nil
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	const op = "service.GetAllOrders"

	orders, err := s.OrderRepo.GetAllOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return orders, nil
}

func (s *OrderService) GetOrderById(ctx context.Context, orderId string) (models.Order, error) {
	const op = "service.GetOrderById"

	order, err := s.OrderRepo.GetOrderById(ctx, orderId)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", op, err)
	}

	return order, nil
}

func (s *OrderService) UpdateOrderById(ctx context.Context, orderId string, order models.Order) error {
	const op = "service.UpdateOrderById"

	err := s.OrderRepo.UpdateOrderById(ctx, orderId, order)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *OrderService) DeleteOrderById(ctx context.Context, orderId string) error {
	const op = "service.DeleteOrderById"

	err := s.OrderRepo.DeleteOrderById(ctx, orderId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *OrderService) CloseOrderById(ctx context.Context, orderId string) error {
	const op = "service.CloseOrderById"

	order, err := s.OrderRepo.GetOrderById(ctx, orderId)
	if err != nil {
		return fmt.Errorf("%s: order not found: %s, %w", op, orderId, err)
	}

	if order.Status == "closed" {
		return fmt.Errorf("%s: order already closed: %s", op, orderId)
	}

	for _, item := range order.Items {
		menuItem, err := s.MenuService.GetMenuItemById(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("%s: %s, %w", op, item.ProductID, err)
		}

		for _, ingredient := range menuItem.Ingredients {
			if !s.InventoryService.HasSufficientStock(ctx, ingredient.IngredientID, ingredient.Quantity*float64(item.Quantity)) {
				return fmt.Errorf("%s: insufficient stock for ingredient: %s", op, ingredient.IngredientID)
			}
		}
	}

	for _, item := range order.Items {
		menuItem, _ := s.MenuService.GetMenuItemById(ctx, item.ProductID)
		for _, ingredient := range menuItem.Ingredients {
			err = s.InventoryService.DeductStock(ctx, ingredient.IngredientID, ingredient.Quantity*float64(item.Quantity))
			if err != nil {
				return fmt.Errorf("%s: failed to deduct stock for ingredient: %s, %w", op, ingredient.IngredientID, err)
			}
		}
	}

	order.Status = "closed"
	err = s.OrderRepo.UpdateOrderById(ctx, orderId, order)
	if err != nil {
		return fmt.Errorf("%s: failed to update order status, %w", op, err)
	}

	return nil
}
