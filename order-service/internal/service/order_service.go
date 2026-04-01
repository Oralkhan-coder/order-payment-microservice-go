package service

import (
	"context"
	"errors"

	"github.com/Oralkhan-coder/order-service/internal/model"
	"github.com/google/uuid"
)

type OrderService struct {
	repo OrderRepo
}

func NewOrderService(repo OrderRepo) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, customerID string, itemName string, amount int64) (string, error) {
	order := model.Order{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		ItemName:   itemName,
		Amount:     amount,
		Status:     model.OrderStatusPending,
	}
	err := s.repo.Create(ctx, order)
	return order.ID, err
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) CancelOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	} else if order.Status != model.OrderStatusPending {
		return errors.New("order is not in pending state")
	}
	return s.repo.UpdateStatus(ctx, id, model.OrderStatusCancelled)
}