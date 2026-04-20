package handler

import (
	"context"
	"time"

	paymentv1 "github.com/Oralkhan-coder/order-payment-proto-generation/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCPaymentHandler struct {
	paymentv1.UnimplementedPaymentServiceServer
	srv PaymentSrv
}

func NewGRPCPaymentHandler(srv PaymentSrv) *GRPCPaymentHandler {
	return &GRPCPaymentHandler{srv: srv}
}

func (h *GRPCPaymentHandler) ProcessPayment(ctx context.Context, req *paymentv1.PaymentRequest) (*paymentv1.PaymentResponse, error) {
	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be positive")
	}

	payment, err := h.srv.ProcessPayment(ctx, req.OrderId, req.Amount)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &paymentv1.PaymentResponse{
		TransactionId: payment.TransactionID,
		Status:        string(payment.Status),
		ProcessedAt:   timestamppb.New(time.Now()),
	}, nil
}
