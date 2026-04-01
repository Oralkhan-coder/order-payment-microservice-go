package transport

import (
	"context"

	"order-service/internal/model"
)

type OrderSrv interface {
	CreateOrder(ctx context.Context, customerID string, itemName string, amount int64) (int, error)
	GetOrder(ctx context.Context, id int) (*model.Order, error)
	CancelOrder(ctx context.Context, id int) error
}