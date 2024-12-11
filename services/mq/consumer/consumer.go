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
	// start background consumer clients
	consumerConn, err := mq.ConnectRabbitMQ("customers")
	if err != nil {
		log.Fatalf("Failed to create rabbitmq connection (consumer): %v", err)
	}
	defer consumerConn.Close()
	consumerClient, err := mq.NewRabbitMQClient(consumerConn)
	if err != nil {
		log.Fatalf("Failed to create rabbitmq client (consumer): %v", err)
	}
	defer consumerClient.Close()
	fmt.Println("Running consumer client...", consumerClient)

	responseConn, err := mq.ConnectRabbitMQ("customers")
	if err != nil {
		log.Fatalf("Failed to create rabbitmq connection (response): %v", err)
	}
	defer responseConn.Close()
	responseClient, err := mq.NewRabbitMQClient(responseConn)
	if err != nil {
		log.Fatalf("Failed to create rabbitmq client (response): %v", err)
	}
	defer responseClient.Close()
	fmt.Println("Running response client...", responseClient)

	blocking := make(chan bool)
	<-blocking
}
