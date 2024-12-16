package proto

import (
	"fmt"
	"os"

	inv "github.com/gonzabosio/transaction/services/proto/inventory"
	order "github.com/gonzabosio/transaction/services/proto/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Services struct {
	Inventory inv.InventoryServiceClient
	Order     order.OrderServiceClient
}

func InitForegroundServices() (*Services, error) {
	inventoryConn, err := grpc.NewClient(os.Getenv("INVENTORY_HOST")+":"+os.Getenv("INVENTORY_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Inventory service: %v", err)
	}

	orderConn, err := grpc.NewClient(os.Getenv("ORDER_HOST")+":"+os.Getenv("ORDER_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Order service: %v", err)
	}

	return &Services{
		Inventory: inv.NewInventoryServiceClient(inventoryConn),
		Order:     order.NewOrderServiceClient(orderConn),
	}, nil
}
