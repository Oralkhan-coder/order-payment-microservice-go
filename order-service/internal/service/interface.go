package service

import (
	"context"

	"github.com/Oralkhan-coder/order-service/internal/model"
)

type OrderRepo interface {
	Create(ctx context.Context, o model.Order) error
	GetByID(ctx context.Context, id string) (*model.Order, error)
	UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error
	GetOrderStatus(ctx context.Context, id string) (*model.OrderStatus, error)
}

type PaymentClient interface {
	AuthorizePayment(ctx context.Context, orderID string, amount int64) (string, error)
}
