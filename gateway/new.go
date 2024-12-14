package gateway

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/gonzabosio/transaction/cache"
	"github.com/gonzabosio/transaction/services/mq"
	"github.com/gonzabosio/transaction/services/proto"
)

type Gateway struct {
	mq         *mq.RabbitClient
	svs        *proto.Services
	cache      *cache.CacheClient
	clientAuth string
}

func NewAPIGateway() (*Gateway, error) {
	publisherConn, err := mq.ConnectRabbitMQ("payments")
	if err != nil {
		return nil, fmt.Errorf("Failed to create rabbitmq connection: %v", err)
	}
	publisherClient, err := mq.NewRabbitMQClient(publisherConn)
	if err != nil {
		return nil, fmt.Errorf("Failed to create rabbitmq client: %v", err)
	}
	svs, err := proto.InitForegroundServices()
	if err != nil {
		return nil, err
	}
	clientAuth := base64.StdEncoding.EncodeToString([]byte(os.Getenv("CLIENT_ID") + ":" + os.Getenv("CLIENT_SECRET")))

	mc := cache.NewMemCachedStorage()

	return &Gateway{
		mq:         &publisherClient,
		svs:        svs,
		cache:      mc,
		clientAuth: clientAuth,
	}, nil
}

func (gw *Gateway) CloseRabbitPublisherChannel() error {
	if err := gw.mq.Close(); err != nil {
		return err
	}
	return nil
}
