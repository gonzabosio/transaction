package gateway

import (
	"fmt"
	"os"

	"github.com/gonzabosio/transaction/services/async/client"
	"github.com/gonzabosio/transaction/services/proto"
)

type gateway struct {
	mq  *client.RabbitClient
	svs *proto.Services
}

func NewAPIGateway() (*gateway, error) {
	conn, err := client.ConnectRabbitMQ(
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		"customers",
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create rabbitmq connection: %v", err)
	}
	client, err := client.NewRabbitMQClient(conn)
	if err != nil {
		return nil, fmt.Errorf("Failed to create rabbitmq client: %v", err)
	}
	svs, err := proto.InitServices()
	if err != nil {
		return nil, err
	}
	return &gateway{
		mq:  &client,
		svs: svs,
	}, nil
}
