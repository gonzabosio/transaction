package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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

func (o *OrderService) NewAccessToken(ctx context.Context, paypalClient *pb.Client) (*pb.AccessToken, error) {
	reqBody := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	req, err := http.NewRequest("POST", "https://www.sandbox.paypal.com/v1/oauth2/token", reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+paypalClient.ClientAuth)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	responseBody := make(map[string]interface{})
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return nil, err
	}
	return &pb.AccessToken{Value: responseBody["access_token"].(string)}, nil
}

func (o *OrderService) NewOrder(ctx context.Context, orderData *pb.Order) (*pb.Result, error) {
	order := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"amount": map[string]string{
					"currency_code": orderData.Currency,
					"value":         orderData.Amount,
				},
			},
		},
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://www.sandbox.paypal.com/v2/checkout/orders", bytes.NewBuffer(orderJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+orderData.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create payment: %v", resp.Status)
	}

	var response map[string]interface{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(body, &response)
	return &pb.Result{Message: "Order created successfully", OrderId: response["id"].(string)}, nil
}

func (o *OrderService) GetOrderDetails(ctx context.Context, order *pb.OrderDetailsRequest) (*pb.OrderDetails, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api-m.sandbox.paypal.com/v2/checkout/orders/%s", order.Id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+order.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	purchaseMap := responseBody["purchase_units"].([]interface{})
	purchaseData := purchaseMap[0].(map[string]interface{})
	amountData := purchaseData["amount"].(map[string]interface{})
	purchase := &pb.PurchaseUnits{}
	purchase.Amount = amountData["value"].(string)
	purchase.Currency = amountData["currency_code"].(string)
	payeeData := purchaseData["payee"].(map[string]interface{})
	purchase.PayeeEmail = payeeData["email_address"].(string)
	purchase.MerchantId = payeeData["merchant_id"].(string)

	return &pb.OrderDetails{Status: responseBody["status"].(string), Purchase: purchase}, nil
}
