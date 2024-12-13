package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gonzabosio/transaction/services/mq"
	"github.com/gonzabosio/transaction/services/proto/inventory/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env: %v", err)
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
	fmt.Println("Running callback client...", callbackClient)

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
			fmt.Printf("Payment callback: %s\nBody: %s\n", message.CorrelationId, string(message.Body))

			callbackBody := make(map[string]interface{})
			json.Unmarshal(message.Body, &callbackBody)
			if v, ok := callbackBody["error_info"]; ok {
				fmt.Printf("%s: %s\n", callbackBody["message"].(string), v)
			} else {
				fmt.Printf("%s\n", callbackBody["message"])
				pId := callbackBody["payment_details"].(map[string]interface{})["productId"].(string)
				productId, err := strconv.ParseInt(pId, 10, 64)
				if err != nil {
					log.Printf("could not parse product id to int64")
				}
				_, err = dbI.Exec("UPDATE product SET stock = stock - 1 WHERE id=$1", productId)
				if err != nil {
					fmt.Printf("Failed to update product %v stock: %v\n", productId, err)
				} else {
					fmt.Printf("Stock modified for the product %v\n", productId)
				}
				// Send email to payer
			}
		}
	}()

	blocking := make(chan bool)
	<-blocking
	callbackClient.Close()
}
