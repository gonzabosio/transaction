package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gonzabosio/transaction/gateway"
	inventory "github.com/gonzabosio/transaction/services/proto/inventory/handlers"
	order "github.com/gonzabosio/transaction/services/proto/order/handlers"
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
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/order", gw.OrderGateway)
	srvPort := os.Getenv("SERVER_PORT")

	go func() {
		log.Printf("API Gateway listening on %s\n", srvPort)
		log.Fatal(http.ListenAndServe(srvPort, r))
	}()

	go inventory.StartInventoryServiceServer()
	go order.StartOrderServiceServer()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
