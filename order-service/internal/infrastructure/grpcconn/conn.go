package grpcconn

import (
	"context"
	"log"
	"os"

	paymentv1 "github.com/Oralkhan-coder/order-payment-proto-generation/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCPaymentConn struct {
	client paymentv1.PaymentServiceClient
}

func NewGRPCPaymentConn() (*GRPCPaymentConn, error) {
	port := os.Getenv("PAYMENT_SERVICE_PORT")
	if port == "" {
		port = "9091"
	}
	log.Printf("Connecting to payment grpc server at port %s", port)
	conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GRPCPaymentConn{client: paymentv1.NewPaymentServiceClient(conn)}, nil
}

func (conn *GRPCPaymentConn) Pay(ctx context.Context, orderID string, amount int64) (string, error) {
	log.Printf("Calling payment gRPC for order %s amount %d", orderID, amount)
	resp, err := conn.client.ProcessPayment(ctx, &paymentv1.PaymentRequest{
		OrderId: orderID,
		Amount:  amount,
	})
	if err != nil {
		log.Printf("Payment gRPC error: %v", err)
		return "", err
	}
	log.Printf("Payment gRPC response: %s", resp.Status)
	return resp.Status, nil
}
