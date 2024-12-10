package proto

import (
	"fmt"
	"os"

	inv "github.com/gonzabosio/transaction/services/proto/inventory"
	order "github.com/gonzabosio/transaction/services/proto/order"
	payment "github.com/gonzabosio/transaction/services/proto/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Services struct {
	Inventory inv.InventoryServiceClient
	Order     order.OrderServiceClient
	Payment   payment.PaymentServiceClient
}

func InitServices() (*Services, error) {
	inventoryConn, err := grpc.NewClient(os.Getenv("INVENTORY_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Inventory service: %v", err)
	}

	orderConn, err := grpc.NewClient(os.Getenv("ORDER_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Order service: %v", err)
	}

	paymentConn, err := grpc.NewClient(os.Getenv("PAYMENT_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Order service: %v", err)
	}

	return &Services{
		Inventory: inv.NewInventoryServiceClient(inventoryConn),
		Order:     order.NewOrderServiceClient(orderConn),
		Payment:   payment.NewPaymentServiceClient(paymentConn),
	}, nil
}
