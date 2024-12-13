package main

import (
	"fmt"
	"log"

	"github.com/gonzabosio/transaction/services/mq"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env: %v", err)
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
			// Send email to payer
		}
	}()

	blocking := make(chan bool)
	<-blocking
	callbackClient.Close()
}
