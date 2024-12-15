package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gonzabosio/transaction/gateway"
	"github.com/gonzabosio/transaction/router"
	email "github.com/gonzabosio/transaction/services/proto/email/handlers"
	inventory "github.com/gonzabosio/transaction/services/proto/inventory/handlers"
	order "github.com/gonzabosio/transaction/services/proto/order/handlers"
	payment "github.com/gonzabosio/transaction/services/proto/payment/handlers"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env: %v", err)
	}

	gw, err := gateway.NewAPIGateway()
	if err != nil {
		log.Fatalf("Failed to create api gateway: %v", err)
	}
	r := router.NewRouter(gw)
	srvPort := os.Getenv("SERVER_PORT")
	go func() {
		log.Printf("API Gateway listening on %s\n", srvPort)
		log.Fatal(http.ListenAndServe(srvPort, r))
	}()

	go inventory.StartInventoryServiceServer()
	go order.StartOrderServiceServer()
	go payment.StartPaymentServiceServer()
	go email.StartEmailServiceServer()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	if err := gw.CloseRabbitPublisherChannel(); err != nil {
		log.Fatalf("Failed to close publisher channel: %s", err)
	} else {
		log.Println("Publisher channel was closed")
	}
}
