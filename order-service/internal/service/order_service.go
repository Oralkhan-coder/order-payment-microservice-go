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
	cache         OrderCache
}

func NewOrderService(repo OrderRepo, paymentClient PaymentGRPCClient, cache OrderCache) *OrderService {
	return &OrderService{
		repo:          repo,
		paymentClient: paymentClient,
		cache:         cache,
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
		_ = s.cache.Delete(ctx, order.ID)
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
	// Invalidate cache after status change
	_ = s.cache.Delete(ctx, order.ID)

	return order.ID, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	// Cache-aside read path
	if cached, err := s.cache.Get(ctx, id); err != nil {
		log.Printf("cache get error for order %s: %v", id, err)
	} else if cached != nil {
		return cached, nil
	}

	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, order)
	return order, nil
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

	if err := s.repo.UpdateStatus(ctx, id, model.OrderStatusCancelled); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, id)
	return nil
}
