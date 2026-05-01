package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/Oralkhan-coder/payment-service/config"
	"github.com/Oralkhan-coder/payment-service/internal/infrastructure/postgres"
	"github.com/Oralkhan-coder/payment-service/internal/infrastructure/rabbitmq"
	"github.com/Oralkhan-coder/payment-service/internal/repository"
	"github.com/Oralkhan-coder/payment-service/internal/service"
	"github.com/Oralkhan-coder/payment-service/internal/transport/grpc"
	"github.com/Oralkhan-coder/payment-service/internal/transport/http"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env", "../.env"); err != nil {
		log.Printf("warning: .env file not loaded: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg := config.InitConfig()

	db, err := infrastructure_postgres.NewConnectionPool(ctx, infrastructure_postgres.NewConfigMust())
	if err != nil {
		fmt.Println("failed to connect to database: ", err)
	}
	defer db.Close()

	rmqConn, err := rabbitmq.NewConnection(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("unable to connect to rabbitmq: %v", err)
	}
	defer rmqConn.Close()

	publisher, err := rabbitmq.NewPublisher(rmqConn.Channel)
	if err != nil {
		log.Fatalf("unable to create rabbitmq publisher: %v", err)
	}

	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo, publisher)

	grpcServer := grpc.NewPaymentGRPCServer(paymentService)
	server := http.NewServer(paymentService)

	log.Printf("starting payment HTTP server on :%s", cfg.HTTPPort)

	go func() {
		grpcServer.Run(ctx)
	}()
	go func() {
		if err := server.Run(ctx, ":"+cfg.HTTPPort); err != nil {
			log.Printf("http server stopped: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Println("payment service shutting down")
}
