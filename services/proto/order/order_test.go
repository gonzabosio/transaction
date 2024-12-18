package order_test

import (
	"context"
	"encoding/base64"
	"os"
	"testing"

	pb "github.com/gonzabosio/transaction/services/proto/order"
	"github.com/gonzabosio/transaction/services/proto/order/handlers"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestOrderServiceFullProcess(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatal(err)
	}
	s := &handlers.OrderService{}
	// sandbox app
	clientAuth := base64.StdEncoding.EncodeToString([]byte(os.Getenv("CLIENT_ID") + ":" + os.Getenv("CLIENT_SECRET")))
	at, err := s.NewAccessToken(context.Background(), &pb.Client{ClientAuth: clientAuth})
	if err != nil {
		t.Fatal(err)
	}
	if at.Value == "" {
		t.Fatal("failed to create access token")
	}

	res, err := s.NewOrder(context.Background(), &pb.Order{AccessToken: at.Value, Currency: "USD", Amount: "30"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Order created successfully", res.Message)

	det, err := s.GetOrderDetails(context.Background(), &pb.OrderDetailsRequest{Id: res.OrderId, AccessToken: at.Value})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "CREATED", det.Status)
}
