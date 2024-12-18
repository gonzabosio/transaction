package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gonzabosio/transaction/services/mq"
	"github.com/gonzabosio/transaction/services/proto/payment"
	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type paymentClient struct {
	client payment.PaymentServiceClient
}

type chanResult struct {
	result *payment.Result
	err    error
}

type paymentRequest struct {
	OrderId     string `json:"order_id"`
	AccessToken string `json:"access_token"`
	ProductId   int64  `json:"product_id"`
}

func NewPaymentServiceClient() (*paymentClient, error) {
	paymentConn, err := grpc.NewClient(os.Getenv("PAYMENT_HOST")+":"+os.Getenv("PAYMENT_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Order service: %v", err)
	}
	return &paymentClient{client: payment.NewPaymentServiceClient(paymentConn)}, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading env: %v\n", err)
	}

	p, err := NewPaymentServiceClient()
	if err != nil {
		log.Fatalf("Failed to create new payment service client: %v", err)
	}

	paymentsConn, err := mq.ConnectRabbitMQ("payments")
	if err != nil {
		log.Fatalf("Failed to create rabbitmq connection (consumer): %v", err)
	}
	defer paymentsConn.Close()
	// start background consumer clients
	consumerClient, err := mq.NewRabbitMQClient(paymentsConn)
	if err != nil {
		log.Fatalf("Failed to create rabbitmq client (consumer): %v", err)
	}
	defer consumerClient.Close()
	fmt.Println("Running consumer client")
	responseClient, err := mq.NewRabbitMQClient(paymentsConn)
	if err != nil {
		log.Fatalf("Failed to create rabbitmq client (response): %v", err)
	}
	defer responseClient.Close()
	fmt.Println("Running response client")

	queue, err := consumerClient.CreateQueue(mq.Q1, true, true)
	if err != nil {
		log.Fatalf("Failed to create queue: %v", err)
	}

	if err := consumerClient.CreateBinding(queue.Name, "", "payment_events"); err != nil {
		log.Fatalf("Failed to create binding: %v", err)
	}

	messageBus, err := consumerClient.Consume(queue.Name, "payment_consumer", false)
	if err != nil {
		log.Fatal("Failed to consume message")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)

	// hard limit on the server
	if err := consumerClient.ApplyQos(10, 0, true); err != nil {
		panic(err)
	}

	go func() {
		for msg := range messageBus {
			r := &chanResult{}
			var productId int64
			retryCount := 2
			g.Go(func() error {
				correlationId := msg.CorrelationId
				log.Printf("New message(%s): %v\n", correlationId, string(msg.Body))
				payload := &paymentRequest{}
				if err := json.Unmarshal(msg.Body, &payload); err != nil {
					log.Printf("Unmarshal payment payload failed: %v\n", err)
					return err
				}
				productId = payload.ProductId
				checkoutCh := make(chan *chanResult)
				for attempt := 1; attempt <= retryCount; attempt++ {
					go func() {
						ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
						defer cancel()
						res, err := p.client.CheckoutOrder(ctx, &payment.CheckoutRequest{OrderId: payload.OrderId, AccessToken: payload.AccessToken})
						if err != nil {
							checkoutCh <- &chanResult{result: nil, err: err}
						} else {
							checkoutCh <- &chanResult{result: res, err: nil}
						}
					}()
					r = <-checkoutCh
					if r.err == nil {
						break
					}
					if attempt < retryCount {
						log.Printf("retrying checkout order; message correlation id: %s\n", correlationId)
						time.Sleep(2 * time.Second)
					} else {
						break
					}
				}
				if err := msg.Ack(false); err != nil {
					return err
				}
				return nil
			})
			if err := g.Wait(); err != nil {
				if err := responseClient.Send(ctx, "payment_callbacks", msg.ReplyTo, amqp091.Publishing{
					ContentType:  "application/json",
					DeliveryMode: amqp091.Persistent,
					Body: []byte(fmt.Sprintf(`{
						"message": "Payment failed before processing it",
						"error_info": "%v"
					}`, err)),
					CorrelationId: msg.CorrelationId,
				}); err != nil {
					log.Printf("Failed to send callback message to producer: %v\n", err)
				}
			} else if r.err != nil {
				if err := responseClient.Send(ctx, "payment_callbacks", msg.ReplyTo, amqp091.Publishing{
					ContentType:  "application/json",
					DeliveryMode: amqp091.Persistent,
					Body: []byte(fmt.Sprintf(`{
						"message": "Payment failed during checkout",
						"error_info": "%v"
					}`, r.err)),
					CorrelationId: msg.CorrelationId,
				}); err != nil {
					log.Printf("Failed to send callback message to producer: %v\n", err)
				}
			} else {
				r.result.ProductId = productId
				b, err := protojson.Marshal(r.result)
				if err != nil {
					log.Printf("Failed to marshal result")
				}
				if err := responseClient.Send(ctx, "payment_callbacks", msg.ReplyTo, amqp091.Publishing{
					ContentType:  "application/json",
					DeliveryMode: amqp091.Persistent,
					Body: []byte(fmt.Sprintf(`{
							"message": "Payment completed",
							"payment_details": %s
							}`, string(b))),
					CorrelationId: msg.CorrelationId,
				}); err != nil {
					log.Printf("Failed to send callback message to producer: %v\n", err)
				}
			}
		}
	}()

	blocking := make(chan bool)
	<-blocking
}
