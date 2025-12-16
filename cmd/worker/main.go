package main

import (
	"Gobank/internal/queue/rabbitmq"
	"Gobank/pkg/config"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := config.LoadConfig()
	conn, err := amqp091.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()
	msg, err := ch.Consume(
		"transfer_event",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Worker started. Waiting for events...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for d := range msg {
			var event rabbitmq.TransferEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error decoding event: %v", err)
				continue
			}
			log.Printf("⚡ Processing Transfer Event: ID=%d Amount=%d Currency=%s",
				event.TransferID, event.Amount, event.Currency)
			//send email or anything else ....
		}
	}()
	<-quit
	log.Println("Worker shutting down...")
}
