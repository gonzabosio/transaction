package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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
	res, err := CaptureOrder(req.OrderId, req.AccessToken)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func CaptureOrder(orderID, accessToken string) (*pb.Result, error) {
	url := fmt.Sprintf("https://api.sandbox.paypal.com/v2/checkout/orders/%s/capture", orderID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var responseBody map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	// fmt.Println("JSON:", string(body))
	json.Unmarshal(body, &responseBody)

	if resp.StatusCode != http.StatusCreated {
		err := responseBody["details"].([]interface{})[0].(map[string]interface{})["issue"].(string)
		return nil, fmt.Errorf("Payment could not be completed: %v", err)
	}

	capture := responseBody["purchase_units"].([]interface{})[0].(map[string]interface{})["payments"].(map[string]interface{})["captures"].([]interface{})[0].(map[string]interface{})
	sellerInfo := capture["seller_receivable_breakdown"].(map[string]interface{})
	netAmount := sellerInfo["net_amount"].(map[string]interface{})["value"].(string)
	paypalFee := sellerInfo["paypal_fee"].(map[string]interface{})["value"].(string)
	refundUrl := capture["links"].([]interface{})[1].(map[string]interface{})["href"].(string)

	return &pb.Result{OrderId: orderID, OrderStatus: responseBody["status"].(string), NetAmount: netAmount, PaypalFee: paypalFee, RefundUrl: refundUrl}, nil
}
