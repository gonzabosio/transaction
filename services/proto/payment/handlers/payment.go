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
	ApiBaseUrl string
}

func StartPaymentServiceServer() {
	lis, err := net.Listen("tcp", os.Getenv("PAYMENT_PORT"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, &PaymentService{ApiBaseUrl: "https://api.sandbox.paypal.com/v2"})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (p *PaymentService) CheckoutOrder(ctx context.Context, req *pb.CheckoutRequest) (*pb.Result, error) {
	res, err := CaptureOrder(p.ApiBaseUrl, req.OrderId, req.AccessToken)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func CaptureOrder(baseUrl, orderID, accessToken string) (*pb.Result, error) {
	url := fmt.Sprintf("%s/checkout/orders/%s/capture", baseUrl, orderID)
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
	payerEmail := responseBody["payment_source"].(map[string]interface{})["paypal"].(map[string]interface{})["email_address"].(string)

	return &pb.Result{OrderId: orderID, OrderStatus: responseBody["status"].(string), NetAmount: netAmount, PaypalFee: paypalFee, PayerEmail: payerEmail, RefundUrl: refundUrl}, nil
}
