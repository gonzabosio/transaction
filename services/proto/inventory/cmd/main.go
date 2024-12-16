package main

import (
	"log"
	"net"
	"os"

	pb "github.com/gonzabosio/transaction/services/proto/inventory"
	invdb "github.com/gonzabosio/transaction/services/proto/inventory/db"
	"github.com/gonzabosio/transaction/services/proto/inventory/handlers"
	"github.com/joho/godotenv"

	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("failed to load enviroment variables: %v", err)
	}
	port := os.Getenv("INVENTORY_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	db, err := invdb.NewInventoryDbConn()
	if err != nil {
		log.Fatalf("Connection to database failed: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	db.SetMaxOpenConns(50)
	defer db.Close()
	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, &handlers.InventoryService{DB: db})

	log.Printf("Inventory service listening on port %s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
