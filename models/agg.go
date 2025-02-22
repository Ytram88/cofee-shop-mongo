package models

type PopularItem struct {
	ProductId string `json:"product_id" bson:"_id"`
	Sold      int    `json:"total_quantity" bson:"total_quantity"`
}
