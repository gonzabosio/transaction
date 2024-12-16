package main

import (
	"log"
	"net"
	"os"

	pb "github.com/gonzabosio/transaction/services/proto/order"
	"github.com/gonzabosio/transaction/services/proto/order/handlers"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("failed to load enviroment variables: %v", err)
	}
	port := os.Getenv("ORDER_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, &handlers.OrderService{})

	log.Printf("Order service listening on port %s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
