package handler

import (
	"context"

	"github.com/Oralkhan-coder/order-service/internal/model"
)

type OrderSrv interface {
	CreateOrder(ctx context.Context, customerID string, itemName string, amount int64) (string, error)
	GetOrder(ctx context.Context, id string) (*model.Order, error)
	CancelOrder(ctx context.Context, id string) error
}
