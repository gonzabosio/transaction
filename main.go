package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gonzabosio/transaction/gateway"
	middleprom "github.com/gonzabosio/transaction/prometheus/middleware"
	"github.com/gonzabosio/transaction/router"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading env: %v\n", err)
	}

	m := middleprom.MetricsSetup()
	m.PrometheusInit()

	gw, err := gateway.NewAPIGateway()
	if err != nil {
		log.Fatalf("Failed to create api gateway: %v", err)
	}
	r := router.NewRouter(gw, m)
	srvPort := os.Getenv("SERVER_PORT")
	go func() {
		log.Printf("API Gateway listening on %s\n", srvPort)
		log.Fatal(http.ListenAndServe(":"+srvPort, r))
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	if err := gw.CloseRabbitPublisherChannel(); err != nil {
		log.Fatalf("Failed to close publisher channel: %s", err)
	} else {
		log.Println("Publisher channel was closed")
	}
}
