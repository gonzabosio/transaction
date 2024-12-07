package handlers

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/gonzabosio/transaction/services/proto/inventory"
	"google.golang.org/grpc"
)

type InventoryService struct {
	pb.UnimplementedInventoryServiceServer
}

func StartInventoryServiceServer() {
	lis, err := net.Listen("tcp", os.Getenv("INVENTORY_PORT"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, &InventoryService{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (i *InventoryService) GetProducts(pr *pb.ProductsRequest, srv pb.InventoryService_GetProductsServer) error {
	return nil
}

func (i *InventoryService) GetStock(ctx context.Context, in *pb.ProductRequest) (*pb.Available, error) {
	return nil, nil
}
