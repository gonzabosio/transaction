package proto

import (
	"fmt"

	invdb "github.com/gonzabosio/transaction/services/proto/inventory/db"
	inv "github.com/gonzabosio/transaction/services/proto/inventory/handlers"
	order "github.com/gonzabosio/transaction/services/proto/order/handlers"
)

type Services struct {
	Inventory *inv.InventoryService
	Order     *order.OrderService
}

func InitServices() (*Services, error) {
	db, err := invdb.NewInventoryDbConn()
	if err != nil {
		return nil, fmt.Errorf("Database error: %v", err)
	}
	return &Services{
		Inventory: &inv.InventoryService{DB: db},
		Order:     &order.OrderService{},
	}, nil
}
