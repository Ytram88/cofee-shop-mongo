package repository

import (
	"cofee-shop-mongo/models"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type InventoryRepository struct {
	collection *mongo.Collection
}

func NewInventoryRepository(db *mongo.Database) *InventoryRepository {
	return &InventoryRepository{
		collection: db.Collection("inventory"),
	}
}

func (r *InventoryRepository) CreateInventoryItem(ctx context.Context, item models.InventoryItem) (string, error) {
	const op = "repository.CreateInventoryItem"
	_, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return item.IngredientID, nil
}

func (r *InventoryRepository) GetAllInventoryItems(ctx context.Context) ([]models.InventoryItem, error) {
	const op = "repository.GetAllInventoryItems"
	var items []models.InventoryItem

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item models.InventoryItem
		if err := cursor.Decode(&item); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		items = append(items, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return items, nil
}

func (r *InventoryRepository) GetInventoryItemById(ctx context.Context, id string) (models.InventoryItem, error) {
	const op = "repository.GetInventoryItemById"
	var item models.InventoryItem

	err := r.collection.FindOne(ctx, bson.M{"ingredient_id": id}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.InventoryItem{}, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return models.InventoryItem{}, fmt.Errorf("%s: %w", op, err)
	}
	return item, nil
}

func (r *InventoryRepository) DeleteInventoryItemById(ctx context.Context, id string) error {
	const op = "repository.DeleteInventoryItemById"
	res, err := r.collection.DeleteOne(ctx, bson.M{"ingredient_id": id})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", op, ErrNotFound)
	}
	return nil
}

func (r *InventoryRepository) UpdateInventoryItemById(ctx context.Context, id string, item models.InventoryItem) error {
	const op = "repository.UpdateInventoryItemById"
	filter := bson.M{"ingredient_id": id}
	update := bson.M{"$set": bson.M{
		"name":     item.Name,
		"quantity": item.Quantity,
		"unit":     item.Unit,
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
