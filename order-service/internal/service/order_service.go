package service

import (
	"context"
	"errors"
	"log"

	"github.com/Oralkhan-coder/order-service/internal/model"
	"github.com/google/uuid"
)

var (
	ErrInvalidAmount         = errors.New("amount must be greater than zero")
	ErrPaymentUnavailable    = errors.New("payment service unavailable")
	ErrOrderAlreadyPaid      = errors.New("paid orders cannot be cancelled")
	ErrOrderAlreadyCancelled = errors.New("order is already cancelled")
	ErrOrderNotFound         = errors.New("order not found")
)

type OrderService struct {
	repo          OrderRepo
	paymentClient PaymentGRPCClient
}

func NewOrderService(repo OrderRepo, paymentClient PaymentGRPCClient) *OrderService {
	return &OrderService{
		repo:          repo,
		paymentClient: paymentClient,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, customerID string, itemName string, amount int64) (string, error) {
	if amount <= 0 {
		return "", ErrInvalidAmount
	}

	order := model.Order{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		ItemName:   itemName,
		Amount:     amount,
		Status:     model.OrderStatusPending,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return "", err
	}

	status, err := s.paymentClient.Pay(ctx, order.ID, order.Amount)
	if err != nil {
		log.Printf("Payment authorization failed for order %s: %v", order.ID, err)
		_ = s.repo.UpdateStatus(ctx, order.ID, model.OrderStatusFailed)
		return order.ID, ErrPaymentUnavailable
	}

	if status == "Authorized" {
		if err := s.repo.UpdateStatus(ctx, order.ID, model.OrderStatusPaid); err != nil {
			return order.ID, err
		}
	} else {
		if err := s.repo.UpdateStatus(ctx, order.ID, model.OrderStatusFailed); err != nil {
			return order.ID, err
		}
	}

	return order.ID, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) GetOrderStatus(ctx context.Context, id string) (*model.OrderStatus, error) {
	return s.repo.GetOrderStatus(ctx, id)
}

func (s *OrderService) CancelOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status == model.OrderStatusPaid {
		return ErrOrderAlreadyPaid
	}
	if order.Status == model.OrderStatusCancelled {
		return ErrOrderAlreadyCancelled
	}
	if order.Status != model.OrderStatusPending {
		return errors.New("only pending orders can be cancelled")
	}

	return s.repo.UpdateStatus(ctx, id, model.OrderStatusCancelled)
}
