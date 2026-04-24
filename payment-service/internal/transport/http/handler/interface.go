package handler

import (
	"context"

	paymentv1 "github.com/Oralkhan-coder/order-payment-proto-generation/payment/v1"
	"github.com/Oralkhan-coder/payment-service/internal/model"
)

type PaymentSrv interface {
	ProcessPayment(ctx context.Context, orderID string, amount int64) (model.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID string) (*model.Payment, error)
	ListPayments(ctx context.Context, req *paymentv1.ListPaymentsRequest) ([]*model.Payment, error)
}
