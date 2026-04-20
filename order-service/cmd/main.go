package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/Oralkhan-coder/order-service/config"
	"github.com/Oralkhan-coder/order-service/internal/infrastructure/grpcconn"
	"github.com/Oralkhan-coder/order-service/internal/infrastructure/postgres"
	"github.com/Oralkhan-coder/order-service/internal/repository"
	"github.com/Oralkhan-coder/order-service/internal/service"
	"github.com/Oralkhan-coder/order-service/internal/transport/grpc"
	"github.com/Oralkhan-coder/order-service/internal/transport/http"
	"github.com/Oralkhan-coder/order-service/pkg"
)

func main() {
	ctx := context.Background()
	cfg := config.InitConfig()
	if err := pkg.RunMigrations(*cfg.Db); err != nil {
		log.Printf("failed to run migrations: %v", err)
	}

	db, err := postgres.NewDB(ctx, *cfg.Db)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer db.Pool.Close()

	paymentClient, err := grpcconn.NewGRPCPaymentConn()
	if err != nil {
		log.Fatalf("unable to connect to payment service: %v", err)
	}

	orderRepo := repository.NewOrderRepository(db.Pool)
	orderService := service.NewOrderService(orderRepo, paymentClient)

	grpcServer := grpc.NewOrderGRPCServer(orderService)
	server := http.NewServer(orderService)

	log.Println("starting the server on :8080")

	grpcServer.Run(ctx)
	server.Run(ctx)
}
