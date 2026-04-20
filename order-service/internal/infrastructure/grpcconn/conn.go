package grpcconn

import (
	"context"
	"log"

	paymentv1 "github.com/Oralkhan-coder/order-payment-proto-generation/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCPaymentConn struct {
	client paymentv1.PaymentServiceClient
}

func NewGRPCPaymentConn(host, port string) (*GRPCPaymentConn, error) {
	address := host + ":" + port
	log.Printf("Connecting to payment grpc server at %s", address)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
