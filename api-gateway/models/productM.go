package models

type ProductResponseSt struct {
	ID        int32   `json:"id"`
	Name      string  `json:"name"`
	Quantity  int32   `json:"quantity"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
}
