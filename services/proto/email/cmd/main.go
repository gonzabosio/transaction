package main

import (
	"log"
	"net"
	"os"

	pb "github.com/gonzabosio/transaction/services/proto/email"
	"github.com/gonzabosio/transaction/services/proto/email/handlers"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("failed to load enviroment variables: %v", err)
	}
	port := os.Getenv("EMAIL_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Printf("Failed to listen: %v\n", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterEmailServiceServer(grpcServer, &handlers.EmailService{})

	log.Printf("Email service listening on port %s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
