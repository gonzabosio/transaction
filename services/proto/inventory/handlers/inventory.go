package handlers

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	"github.com/gonzabosio/transaction/model"
	pb "github.com/gonzabosio/transaction/services/proto/inventory"
	invdb "github.com/gonzabosio/transaction/services/proto/inventory/db"
	"google.golang.org/grpc"
)

type InventoryService struct {
	pb.UnimplementedInventoryServiceServer
	DB *sql.DB
}

func StartInventoryServiceServer() {
	lis, err := net.Listen("tcp", os.Getenv("INVENTORY_PORT"))
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
	pb.RegisterInventoryServiceServer(grpcServer, &InventoryService{DB: db})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (i *InventoryService) GetProducts(pr *pb.ProductsRequest, srv pb.InventoryService_GetProductsServer) error {
	return nil
}

func (i *InventoryService) GetStock(ctx context.Context, payload *pb.ProductRequest) (*pb.Available, error) {
	row := i.DB.QueryRow(`SELECT * FROM product WHERE id=$1`, payload.ProductId)
	product := &model.Product{}
	if err := row.Scan(&product.Id, &product.Name, &product.Stock, &product.Price); err != nil {
		return nil, err
	}
	result := &pb.Available{Stock: product.Stock, Price: product.Price, Name: product.Name}
	result.IsAvailable = product.Stock > 0
	if result.IsAvailable {
		_, err := i.DB.Exec("UPDATE product SET stock = stock - 1 WHERE id=$1", product.Id)
		if err != nil {
			return nil, err
		} else {
			result.Stock -= 1
		}
	}
	return result, nil
}
