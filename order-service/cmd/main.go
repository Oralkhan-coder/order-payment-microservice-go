package main

import (
	"context"
	"log"

	"github.com/Oralkhan-coder/order-service/config"
	"github.com/Oralkhan-coder/order-service/internal/infrastructure/grpcconn"
	"github.com/Oralkhan-coder/order-service/internal/infrastructure/postgres"
	"github.com/Oralkhan-coder/order-service/internal/repository"
	"github.com/Oralkhan-coder/order-service/internal/service"
	"github.com/Oralkhan-coder/order-service/internal/transport/grpc"
	"github.com/Oralkhan-coder/order-service/internal/transport/http"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env", "../.env"); err != nil {
		log.Printf("warning: .env file not loaded: %v", err)
	}

	ctx := context.Background()
	cfg := config.InitConfig()

	db, err := postgres.NewConnectionPool(ctx, postgres.NewConfigMust())
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer db.Close()

	paymentClient, err := grpcconn.NewGRPCPaymentConn(cfg.PaymentServiceHost, cfg.PaymentServicePort)
	if err != nil {
		log.Fatalf("unable to connect to payment service: %v", err)
	}

	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, paymentClient)

	grpcServer := grpc.NewOrderGRPCServer(orderService)
	server := http.NewServer(orderService)

	go func() {
		grpcServer.Run(ctx)
	}()
	log.Println("starting the server on :8080")
	server.Run(ctx)
}
