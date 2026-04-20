package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	paymentv1 "github.com/Oralkhan-coder/order-payment-proto-generation/payment/v1"
	"github.com/Oralkhan-coder/payment-service/internal/transport/grpc/handler"
	"google.golang.org/grpc"
)

type PaymentGRPCServer struct {
	grpcServer     *grpc.Server
	paymentHandler *handler.GRPCPaymentHandler
}

func NewPaymentGRPCServer(srv handler.PaymentSrv) *PaymentGRPCServer {
	grpcPaymentHandler := handler.NewGRPCPaymentHandler(srv)
	grpcServer := grpc.NewServer()
	paymentv1.RegisterPaymentServiceServer(grpcServer, grpcPaymentHandler)

	return &PaymentGRPCServer{grpcServer, grpcPaymentHandler}
}

func (s *PaymentGRPCServer) Run(ctx context.Context) {
	grpcOrderServicePort := os.Getenv("GRPC_PAYMENT_SERVICE_PORT")
	if grpcOrderServicePort == "" {
		grpcOrderServicePort = "9091"
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
