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
	callbackConn, err := mq.ConnectRabbitMQ("customers")
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

	blocking := make(chan bool)
	<-blocking
}
