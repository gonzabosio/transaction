package gateway

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/gonzabosio/transaction/services/async/client"
	"github.com/gonzabosio/transaction/services/proto"
)

type gateway struct {
	mq         *client.RabbitClient
	svs        *proto.Services
	clientAuth string
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
	clientAuth := base64.StdEncoding.EncodeToString([]byte(os.Getenv("CLIENT_ID") + ":" + os.Getenv("CLIENT_SECRET")))

	return &gateway{
		mq:         &client,
		svs:        svs,
		clientAuth: clientAuth,
	}, nil
}
