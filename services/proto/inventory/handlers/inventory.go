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

func (i *InventoryService) GetProducts(req *pb.ProductsRequest, srv pb.InventoryService_GetProductsServer) error {
	rows, err := i.DB.Query(`SELECT * FROM product`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var product pb.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Stock, &product.Price); err != nil {
			return err
		}
		if err := srv.Send(&pb.Product{
			Id:    product.Id,
			Name:  product.Name,
			Stock: product.Stock,
			Price: product.Price,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (i *InventoryService) GetStock(ctx context.Context, payload *pb.ProductRequest) (*pb.Available, error) {
	row := i.DB.QueryRow(`SELECT * FROM product WHERE id=$1`, payload.ProductId)
	product := &model.Product{}
	if err := row.Scan(&product.Id, &product.Name, &product.Stock, &product.Price); err != nil {
		return nil, err
	}
	result := &pb.Available{ProductId: product.Id, Stock: product.Stock, Price: product.Price, Name: product.Name}
	result.IsAvailable = product.Stock > 0
	return result, nil
}
