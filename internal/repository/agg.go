package repository

import (
	"cofee-shop-mongo/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ReportRepository struct {
	db *mongo.Database
}

func NewReportRepository(db *mongo.Database) *ReportRepository {
	return &ReportRepository{db}
}

func (r *ReportRepository) GetTotalSales(ctx context.Context) (float64, error) {
	const op = "repository.GetTotalSales"
	collection := r.db.Collection("orders")

	pipeline := []bson.M{
		{"$unwind": "$items"},
		{
			"$lookup": bson.M{
				"from":         "menu",
				"localField":   "items.product_id",
				"foreignField": "product_id",
				"as":           "product",
			},
		},
		{"$unwind": "$product"},
		{
			"$project": bson.M{
				"total_price": bson.M{
					"$multiply": []interface{}{"$product.price", "$items.quantity"},
				},
			},
		},
		{
			"$group": bson.M{
				"_id":         nil,
				"total_sales": bson.M{"$sum": "$total_price"},
			},
		},
		{"$project": bson.M{"_id": 0}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalSales float64 `bson:"total_sales"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	} else {
		return 0, fmt.Errorf("%s: %w", op, ErrNotFound)
	}

	return result.TotalSales, nil
}

func (r *ReportRepository) GetPopularItems(ctx context.Context) ([]models.PopularItem, error) {
	const op = "repository.GetPopularItems"
	collection := r.db.Collection("orders")

	pipeline := []bson.M{
		{"$unwind": "$items"},
		{
			"$group": bson.M{
				"_id":            "$items.product_id",
				"total_quantity": bson.M{"$sum": "$items.quantity"},
			},
		},
		{"$sort": bson.M{"total_quantity": -1}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	var popularItems []models.PopularItem
	for cursor.Next(ctx) {
		var item models.PopularItem
		if err := cursor.Decode(&item); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		popularItems = append(popularItems, item)
	}

	if len(popularItems) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrNotFound)
	}

	return popularItems, nil
}
