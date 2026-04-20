package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	orderv1 "github.com/Oralkhan-coder/order-payment-proto-generation/order/v1"
	"github.com/Oralkhan-coder/order-service/internal/transport/grpc/handler"
	"google.golang.org/grpc"
)

type OrderGRPCServer struct {
	grpcServer   *grpc.Server
	orderHandler *handler.GRPCOrderHandler
}

func NewOrderGRPCServer(srv handler.OrderSrv) *OrderGRPCServer {
	grpcOrderHandler := handler.NewGRPCOrderHandler(srv)
	grpcServer := grpc.NewServer()
	orderv1.RegisterOrderServiceServer(grpcServer, grpcOrderHandler)

	return &OrderGRPCServer{
		grpcServer:   grpcServer,
		orderHandler: grpcOrderHandler,
	}
}

func (s *OrderGRPCServer) Run(ctx context.Context) {
	grpcOrderServicePort := os.Getenv("GRPC_ORDER_SERVICE_PORT")
	if grpcOrderServicePort == "" {
		grpcOrderServicePort = "9090"
	}
	listener, err := net.Listen("tcp", ":"+grpcOrderServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("gRPC server is running on %s\n", grpcOrderServicePort)

	go func() {
		if err := s.grpcServer.Serve(listener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("Shutting down gRPC server...")
	s.grpcServer.GracefulStop()
}
