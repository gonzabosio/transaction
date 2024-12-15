package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gonzabosio/transaction/services/mq"
	"github.com/gonzabosio/transaction/services/proto/email"
	"github.com/gonzabosio/transaction/services/proto/inventory/db"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type emailClient struct {
	client email.EmailServiceClient
}

func NewEmailServiceClient() (*emailClient, error) {
	emailConn, err := grpc.NewClient(os.Getenv("EMAIL_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Order service: %v", err)
	}
	return &emailClient{client: email.NewEmailServiceClient(emailConn)}, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env: %v", err)
	}

	e, err := NewEmailServiceClient()
	if err != nil {
		log.Fatalf("Failed to create email service client: %v", err)
	}

	dbI, err := db.NewInventoryDbConn()
	if err != nil {
		log.Fatalf("Failed to connect to inventory database: %v", err)
	}

	callbackConn, err := mq.ConnectRabbitMQ("payments")
	if err != nil {
		log.Fatalf("Failed to create rabbitmq connection (callback): %v", err)
	}
	defer callbackConn.Close()
	callbackClient, err := mq.NewRabbitMQClient(callbackConn)
	if err != nil {
		log.Fatalf("Failed to create rabbitmq client (callback): %v", err)
	}
	defer callbackClient.Close()
	fmt.Println("Running callback client")

	queue, err := callbackClient.CreateQueue(mq.Q1, true, true)
	if err != nil {
		panic(err)
	}
	if err := callbackClient.CreateBinding(queue.Name, queue.Name, "payment_callbacks"); err != nil {
		panic(err)
	}

	messages, err := callbackClient.Consume(queue.Name, "payment_consumer", true)
	if err != nil {
		panic(err)
	}
	go func() {
		for message := range messages {
			log.Printf("Callback: %s", message.CorrelationId)

			callbackBody := make(map[string]interface{})
			json.Unmarshal(message.Body, &callbackBody)
			if v, ok := callbackBody["error_info"]; ok {
				log.Printf("%s: %s\n", callbackBody["message"].(string), v)
			} else {
				fmt.Printf("%s\n", callbackBody["message"])
				pId := callbackBody["payment_details"].(map[string]interface{})["productId"].(string)
				productId, err := strconv.ParseInt(pId, 10, 64)
				orderId := callbackBody["payment_details"].(map[string]interface{})["orderId"].(string)
				if err != nil {
					log.Printf("could not parse product id to int64")
				}
				_, err = dbI.Exec("UPDATE product SET stock = stock - 1 WHERE id=$1", productId)
				if err != nil {
					log.Printf("Failed to update product %v stock: %v\n", productId, err)
				} else {
					log.Printf("Stock modified for the product %v\n", productId)
				}
				row := dbI.QueryRow(`SELECT name FROM product WHERE id=$1`, productId)
				var prodName string
				if err := row.Scan(&prodName); err != nil {
					fmt.Printf("Failed to scan row: %v", err)
				} else {
					// send email
					ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
					defer cancel()
					res, err := e.client.SendEmail(ctx, &email.EmailRequest{
						Subject:    "PayPal Purchase",
						BodyText:   fmt.Sprintf("Hello! We inform you that your purchase of -%s- was successful 😃.\nIf you want to make a refund go to the link below:\nhttps://myexample.business.com/refund?order_id=%s", prodName, orderId),
						PayerEmail: os.Getenv("PAYER_TEST_EMAIL"),
					})
					if err != nil {
						log.Printf("Email was not sent: %v", err)
					} else {
						log.Println(res.Message)
					}
				}
			}
		}
	}()

	blocking := make(chan bool)
	<-blocking
	callbackClient.Close()
}
