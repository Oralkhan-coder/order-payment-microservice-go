package main

import (
	"context"
	"log"

	"github.com/Oralkhan-coder/payment-service/config"
	"github.com/Oralkhan-coder/payment-service/internal/infrastructure/postgres"
	"github.com/Oralkhan-coder/payment-service/internal/repository"
	"github.com/Oralkhan-coder/payment-service/internal/service"
	"github.com/Oralkhan-coder/payment-service/internal/transport"
	"github.com/Oralkhan-coder/payment-service/pkg"
)

func main() {
	ctx := context.Background()
	cfg := config.InitConfig()

	if err := pkg.RunMigrations(cfg.Db); err != nil {
		log.Printf("failed to run migrations: %v", err)
	}

	db, err := postgres.NewDB(ctx, cfg.Db)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer db.Pool.Close()

	paymentRepo := repository.NewPaymentRepository(db.Pool)
	paymentService := service.NewPaymentService(paymentRepo)
	server := transport.NewServer(paymentService)

	log.Println("starting the server on :8081")

	if err := server.Run(ctx, ":8081"); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
