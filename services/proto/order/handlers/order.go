package handlers

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/gonzabosio/transaction/services/proto/order"
	"google.golang.org/grpc"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
}

func StartOrderServiceServer() {
	lis, err := net.Listen("tcp", os.Getenv("ORDER_PORT"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, &OrderService{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (o *OrderService) NewAccessToken(context.Context, *pb.Client) (*pb.AccessToken, error) {
	return nil, nil
}

func (o *OrderService) NewOrder(context.Context, *pb.Order) (*pb.Result, error) {
	return nil, nil
}

func (o *OrderService) GetOrderDetails(context.Context, *pb.OrderID) (*pb.OrderDetails, error) {
	return nil, nil
}
