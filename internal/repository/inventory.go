package repository

import (
	"cofee-shop-mongo/models"
	"context"
	"errors"

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
	_, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return "", err
	}
	return item.IngredientID, nil
}

func (r *InventoryRepository) GetAllInventoryItems(ctx context.Context) ([]models.InventoryItem, error) {
	var items []models.InventoryItem

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item models.InventoryItem
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *InventoryRepository) GetInventoryItemById(ctx context.Context, id string) (models.InventoryItem, error) {
	var item models.InventoryItem
	err := r.collection.FindOne(ctx, bson.M{"ingredient_id": id}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.InventoryItem{}, errors.New("item not found")
		}
		return models.InventoryItem{}, err
	}
	return item, nil
}

func (r *InventoryRepository) DeleteInventoryItemById(ctx context.Context, id string) error {
	res, err := r.collection.DeleteOne(ctx, bson.M{"ingredient_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		errors.New("item not found")
	}
	return nil
}

func (r *InventoryRepository) UpdateInventoryItemById(ctx context.Context, id string, item models.InventoryItem) error {
	filter := bson.M{"ingredient_id": id}
	update := bson.M{"$set": bson.M{
		"name":     item.Name,
		"quantity": item.Quantity,
		"unit":     item.Unit,
	}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("item not found")
	}
	return nil
}
