package handlers

import (
	"context"
	"database/sql"

	"github.com/gonzabosio/transaction/model"
	pb "github.com/gonzabosio/transaction/services/proto/inventory"
)

type InventoryService struct {
	pb.UnimplementedInventoryServiceServer
	DB *sql.DB
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
