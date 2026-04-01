package transport

import (
	"context"

	"github.com/Oralkhan-coder/payment-service/internal/model"
)

type PaymentSrv interface {
	ProcessPayment(ctx context.Context, orderID string, amount int64) (model.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID string) (*model.Payment, error)
}
