package service

import (
	"context"

	"order-service/internal/model"
)

type OrderService struct {
	repo OrderRepo
}

func NewOrderService(repo OrderRepo) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, customerID string, itemName string, amount int64) (int, error) {
	order := model.Order{
		CustomerID: customerID,
		ItemName:   itemName,
		Amount:     amount,
		Status:     model.OrderStatusPending,
	}
	return s.repo.Create(ctx, order)
}

func (s *OrderService) GetOrder(ctx context.Context, id int) (*model.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) CancelOrder(ctx context.Context, id int) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	} else if order.Status != model.OrderStatusPending {
		return errors.New("order is not in pending state")
	}
	return s.repo.UpdateStatus(ctx, id, model.OrderStatusCancelled)
}	