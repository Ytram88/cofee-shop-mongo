package repository

import (
	"cofee-shop-mongo/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{
		collection: db.Collection("orders"),
	}
}
func (r *OrderRepository) CreateOrder(ctx context.Context, order models.Order) (string, error) {
	const op = "repository.CreateOrder"
	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return order.ProductId, nil
}

func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	const op = "repository.GetAllOrders"
	var orders []models.Order

	cursor, err := r.collection.Find(ctx, bson.D{})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var order models.Order
		if err = cursor.Decode(&order); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderById(ctx context.Context, orderId string) (models.Order, error) {
	const op = "repository.GetOrderById"
	var order models.Order

	err := r.collection.FindOne(ctx, bson.D{{"order_id", orderId}}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Order{}, fmt.Errorf("%s: %w", op, ErrNotFound)

		}
		return models.Order{}, fmt.Errorf("%s: %w", op, err)
	}

	return order, nil
}

func (r *OrderRepository) UpdateOrderById(ctx context.Context, orderId string, order models.Order) error {
	const op = "repository.UpdateOrderById"
	filter := bson.M{"order_id": orderId}
	update := bson.M{"$set": bson.M{
		"customer_name": order.CustomerName,
		"items":         order.Items,
		"status":        order.Status,
		"created_at":    order.CreatedAt,
	}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("%s: %w", op, ErrNotFound)
	}
	return nil
}

func (r *OrderRepository) DeleteOrderById(ctx context.Context, orderId string) error {
	const op = "repository.DeleteOrderById"
	res, err := r.collection.DeleteOne(ctx, bson.D{{"order_id", orderId}})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", op, ErrNotFound)
	}
	return nil
}
