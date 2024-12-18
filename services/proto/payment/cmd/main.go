package main

import (
	"log"
	"net"
	"os"

	pb "github.com/gonzabosio/transaction/services/proto/payment"
	"github.com/gonzabosio/transaction/services/proto/payment/handlers"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("failed to load enviroment variables: %v", err)
	}
	port := os.Getenv("PAYMENT_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, &handlers.PaymentService{ApiBaseUrl: "https://api.sandbox.paypal.com/v2"})

	log.Printf("Payment service listening on port %s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
