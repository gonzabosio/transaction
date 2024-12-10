package handlers

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/gonzabosio/transaction/services/proto/payment"
	"google.golang.org/grpc"
)

type PaymentService struct {
	pb.UnimplementedPaymentServiceServer
}

func StartPaymentServiceServer() {
	lis, err := net.Listen("tcp", os.Getenv("PAYMENT_PORT"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, &PaymentService{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (p *PaymentService) CheckoutOrder(ctx context.Context, req *pb.CheckoutRequest) (*pb.Result, error) {
	return &pb.Result{}, nil
}
