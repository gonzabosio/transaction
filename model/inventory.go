package model

type Product struct {
	Id    int64
	Name  string
	Stock int64
	Price float32
}

type ProductRequest struct {
	ProductId int64 `json:"product_id" validate:"required"`
}
