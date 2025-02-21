package repository

import (
	"cofee-shop-mongo/models"
	"context"
	"errors"
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
	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return "", err
	}
	return order.ProductId, nil
}
func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order

	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var order models.Order
		if err = cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
func (r *OrderRepository) GetOrderById(ctx context.Context, OrderId string) (models.Order, error) {
	var order models.Order

	err := r.collection.FindOne(ctx, bson.D{{"order_id", OrderId}}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Order{}, errors.New("item not found")
		}
		return models.Order{}, err
	}

	return order, nil
}
func (r *OrderRepository) UpdateOrderById(ctx context.Context, OrderId string, order models.Order) error {
	filter := bson.M{"order_id": OrderId}
	update := bson.M{"$set": bson.M{
		"CustomerName": order.CustomerName,
		"Items":        order.Items,
		"Status":       order.Status,
		"CreatedAt":    order.CreatedAt,
	}}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("menu item not found")
	}
	return nil
}
func (r *OrderRepository) DeleteOrderById(ctx context.Context, OrderId string) error {
	res, err := r.collection.DeleteOne(ctx, bson.D{{"order_id", OrderId}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("menu item not found")
	}
	return nil
}
