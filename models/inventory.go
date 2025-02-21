package models

type InventoryItem struct {
	IngredientID string  `bson:"ingredient_id" json:"ingredient_id"`
	Name         string  `bson:"name" json:"name"`
	Quantity     float64 `bson:"quantity" json:"quantity"`
	Unit         string  `bson:"unit" json:"unit"`
}
