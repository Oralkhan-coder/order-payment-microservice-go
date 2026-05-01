package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Oralkhan-coder/notification-service/internal/consumer"
	"github.com/Oralkhan-coder/notification-service/internal/infrastructure"
	"github.com/Oralkhan-coder/notification-service/internal/repository"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	rmq, err := infrastructure.NewRabbitMQ(rabbitURL)
	if err != nil {
		log.Fatalf("failed to initialize rabbitmq: %v", err)
	}
	defer func() {
		_ = rmq.Channel.Close()
		_ = rmq.Conn.Close()
	}()

	store := repository.NewIdempotencyStore()
	c := consumer.NewNotificationConsumer(rmq.Channel, store)
	log.Println("notification consumer started")
	if err := c.Start(ctx); err != nil {
		log.Fatalf("notification consumer stopped with error: %v", err)
	}

	log.Println("notification service shutting down gracefully")
}
