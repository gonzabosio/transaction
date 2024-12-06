package gateway

import (
	"fmt"
	"os"

	"github.com/gonzabosio/transaction/services/async/client"
)

type gateway struct {
	mq *client.RabbitClient
}

func NewAPIGateway() (*gateway, error) {
	conn, err := client.ConnectRabbitMQ(
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		"orders",
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create rabbitmq connection: %v", err)
	}
	client, err := client.NewRabbitMQClient(conn)
	if err != nil {
		return nil, fmt.Errorf("Failed to create rabbitmq client: %v", err)
	}
	return &gateway{
		mq: &client,
	}, nil
}
