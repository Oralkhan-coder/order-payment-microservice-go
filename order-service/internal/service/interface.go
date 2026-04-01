package service

import (
	"context"

	"order-service/internal/model"
)

type OrderRepo interface {
	Create(ctx context.Context, o model.Order) (int, error)
	GetByID(ctx context.Context, id int) (*model.Order, error)
	UpdateStatus(ctx context.Context, id int, status model.OrderStatus) error
}