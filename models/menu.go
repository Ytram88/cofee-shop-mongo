package models

type MenuItem struct {
	ProductId   string               `bson:"product_id" json:"product_id"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	Price       float64              `bson:"price" json:"price"`
	Ingredients []MenuItemIngredient `bson:"ingredients" json:"ingredients"`
}

type MenuItemIngredient struct {
	IngredientID string  `bson:"ingredient_id" json:"ingredient_id"`
	Quantity     float64 `bson:"quantity" json:"quantity"`
}
