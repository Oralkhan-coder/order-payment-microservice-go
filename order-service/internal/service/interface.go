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

type PaymentGRPCClient interface {
	Pay(ctx context.Context, orderID string, amount int64) (string, error)
}

type OrderCache interface {
	Get(ctx context.Context, id string) (*model.Order, error)
	Set(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id string) error
}
