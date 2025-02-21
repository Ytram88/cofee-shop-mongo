package repository

import (
	"cofee-shop-mongo/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MenuRepository struct {
	collection *mongo.Collection
}

func NewMenuRepository(db *mongo.Database) *MenuRepository {
	return &MenuRepository{
		collection: db.Collection("menu"),
	}
}

func (r *MenuRepository) CreateMenuItem(ctx context.Context, item models.MenuItem) (string, error) {
	_, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return "", err
	}
	return item.ProductId, nil
}

func (r *MenuRepository) GetAllMenuItems(ctx context.Context) ([]models.MenuItem, error) {
	var items []models.MenuItem

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item models.MenuItem
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *MenuRepository) GetMenuItemById(ctx context.Context, id string) (models.MenuItem, error) {
	var item models.MenuItem
	err := r.collection.FindOne(ctx, bson.M{"product_id": id}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.MenuItem{}, errors.New("menu item not found")
		}
		return models.MenuItem{}, err
	}
	return item, nil
}

func (r *MenuRepository) UpdateMenuItemById(ctx context.Context, id string, item models.MenuItem) error {
	filter := bson.M{"product_id": id}
	update := bson.M{"$set": bson.M{
		"name":        item.Name,
		"description": item.Description,
		"price":       item.Price,
		"ingredients": item.Ingredients,
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

func (r *MenuRepository) DeleteMenuItemById(ctx context.Context, id string) error {
	res, err := r.collection.DeleteOne(ctx, bson.M{"product_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("menu item not found")
	}
	return nil
}
