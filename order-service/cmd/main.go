package main

import (
	"context"
	"log"
	"net"
	"os"

	orderv1 "github.com/Oralkhan-coder/order-payment-proto-generation/order/v1"
	"github.com/Oralkhan-coder/order-service/config"
	"github.com/Oralkhan-coder/order-service/internal/infrastructure/grpcconn"
	"github.com/Oralkhan-coder/order-service/internal/infrastructure/postgres"
	"github.com/Oralkhan-coder/order-service/internal/repository"
	"github.com/Oralkhan-coder/order-service/internal/service"
	"github.com/Oralkhan-coder/order-service/internal/transport/http"
	"github.com/Oralkhan-coder/order-service/pkg"
	"google.golang.org/grpc"
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

	grpcOrderServicePort := os.Getenv("GRPC_ORDER_SERVICE_PORT")
	if grpcOrderServicePort == "" {
		grpcOrderServicePort = "9090"
	}
	listener, err := net.Listen("tcp", ":"+grpcOrderServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	orderv1.RegisterOrderServiceServer(grpcServer, orderService)

	server := http.NewServer(orderService)

	log.Println("starting the server on :8080")

	server.Run(ctx)
}
