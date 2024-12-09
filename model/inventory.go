package model

type Product struct {
	Id    int
	Name  string
	Stock int64
}

type ProductRequest struct {
	Name string `json:"name" validate:"required"`
}
