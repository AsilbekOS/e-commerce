package models

type OrderResponseJson struct {
	ProductID int32 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}
