package proto

import (
	inv "github.com/gonzabosio/transaction/services/proto/inventory/handlers"
	order "github.com/gonzabosio/transaction/services/proto/order/handlers"
)

type Services struct {
	Inventory *inv.InventoryService
	Order     *order.OrderService
}

func NewServices() *Services {
	return &Services{
		Inventory: &inv.InventoryService{},
		Order:     &order.OrderService{},
	}
}
