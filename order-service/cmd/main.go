package main

import (
	"context"
	"log"

	"github.com/Oralkhan-coder/order-service/internal/http"
	"github.com/Oralkhan-coder/order-service/internal/infrastructure/payment"
	"github.com/Oralkhan-coder/order-service/internal/repository"
	"github.com/Oralkhan-coder/order-service/internal/service"
	"github.com/Oralkhan-coder/order-service/pkg"
)

func main() {
	ctx := context.Background()
	cfg := pkg.Config{
		Database: "done_db",
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "postgres",
	}

	if err := pkg.RunMigrations(cfg); err != nil {
		log.Printf("failed to run migrations: %v", err)
	}

	db, err := pkg.NewDB(ctx, cfg)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer db.Pool.Close()

	// External adapter (Payment client)
	paymentClient := payment.NewPaymentClient("http://localhost:8081")

	// Repository
	orderRepo := repository.NewOrderRepository(db.Pool)

	// Use Case
	orderService := service.NewOrderService(orderRepo, paymentClient)

	// Delivery
	server := http.NewServer(orderService)

	log.Println("starting the server on :8080")

	server.Run(ctx)
}
