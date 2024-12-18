package gateway

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/gonzabosio/transaction/cache"
	"github.com/gonzabosio/transaction/services/mq"
	"github.com/gonzabosio/transaction/services/proto"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type Gateway struct {
	mq         *mq.RabbitClient // publisher
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

func (gw Gateway) RunPaymentCheckoutTasks(orderId, accessToken string, productId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	id := uuid.New()
	if err := gw.mq.Send(ctx, "payment_events", "", amqp091.Publishing{
		ContentType:   "application/json",
		DeliveryMode:  amqp091.Persistent,
		ReplyTo:       mq.Q1,
		CorrelationId: fmt.Sprintf("payment_%s", id.String()),
		Body: []byte(fmt.Sprintf(`{
			"order_id": "%s",
			"access_token": "%s",
			"product_id": %d
		}`, orderId, accessToken, productId)),
	}); err != nil {
		return err
	}
	return nil
}
