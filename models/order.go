package models

type Order struct {
	ProductId    string      `bson:"order_id" json:"order_id"`
	CustomerName string      `bson:"customer_name" json:"customer_name"`
	Items        []OrderItem `bson:"items" json:"items"`
	Status       string      `bson:"status" json:"status"`
	CreatedAt    string      `bson:"created_at" json:"created_at"`
}

type OrderItem struct {
	ProductID string `bson:"product_id" json:"product_id"`
	Quantity  int    `bson:"quantity" json:"quantity"`
}
