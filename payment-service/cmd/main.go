package main

import (
	"context"
	"log"

	"github.com/Oralkhan-coder/payment-service/internal/http"
	"github.com/Oralkhan-coder/payment-service/internal/repository"
	"github.com/Oralkhan-coder/payment-service/internal/service"
	"github.com/Oralkhan-coder/payment-service/pkg"
)

func main() {
	ctx := context.Background()
	cfg := pkg.Config{
		Database: "payment_db",
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

	paymentRepo := repository.NewPaymentRepository(db.Pool)
	paymentService := service.NewPaymentService(paymentRepo)
	server := http.NewServer(paymentService)

	log.Println("starting the server on :8081")

	if err := server.Run(ctx, ":8081"); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
