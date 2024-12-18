package inventory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/gonzabosio/transaction/services/proto/inventory"
	"github.com/gonzabosio/transaction/services/proto/inventory/handlers"
)

func setupService() (*handlers.InventoryService, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open stub database connection: %v", err)
	}
	s := &handlers.InventoryService{DB: db}
	return s, mock, nil
}

func TestInventoryServiceStock(t *testing.T) {
	s, mock, err := setupService()
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectQuery("SELECT \\* FROM product WHERE id=?").WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "stock", "price"}).AddRow(1, "Product 1", 1, 150))

	res, err := s.GetStock(context.Background(), &pb.ProductRequest{ProductId: 1})
	if err != nil {
		t.Fatalf("failed to get stock: %v", err)
	}
	if res.IsAvailable == false {
		t.Fatalf("is available: expected was true but got: %v", res.IsAvailable)
	}
}
func TestInventoryServiceNoStock(t *testing.T) {
	s, mock, err := setupService()
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectQuery("SELECT \\* FROM product WHERE id=?").WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "stock", "price"}).AddRow(1, "Product 1", 0, 150))
	res, err := s.GetStock(context.Background(), &pb.ProductRequest{ProductId: 1})
	if err != nil {
		t.Fatalf("failed to get stock: %v", err)
	}
	if res.IsAvailable == true {
		t.Fatalf("is available: expected was false but got: %v", res.IsAvailable)
	}
}
