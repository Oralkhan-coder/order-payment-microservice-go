package service

import (
	"context"
	"errors"
	"log"

	"github.com/Oralkhan-coder/order-service/internal/model"
	"github.com/google/uuid"
)

var (
	ErrInvalidAmount       = errors.New("amount must be greater than zero")
	ErrPaymentUnavailable  = errors.New("payment service unavailable")
	ErrOrderAlreadyPaid    = errors.New("paid orders cannot be cancelled")
	ErrOrderAlreadyCancelled = errors.New("order is already cancelled")
	ErrOrderNotFound       = errors.New("order not found")
)

type OrderService struct {
	repo          OrderRepo
	paymentClient PaymentClient
}

func NewOrderService(repo OrderRepo, paymentClient PaymentClient) *OrderService {
	return &OrderService{
		repo:          repo,
		paymentClient: paymentClient,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, customerID string, itemName string, amount int64) (string, error) {
	// 0. Validate amount
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

	// 1. Create order in Pending state
	if err := s.repo.Create(ctx, order); err != nil {
		return "", err
	}

	// 2. Call Payment Service
	status, err := s.paymentClient.AuthorizePayment(ctx, order.ID, order.Amount)
	if err != nil {
		log.Printf("Payment authorization failed for order %s: %v", order.ID, err)
		// 3. Mark as Failed if payment service call fails (Requirement 4.3: mark as Failed or remain Pending)
		_ = s.repo.UpdateStatus(ctx, order.ID, model.OrderStatusFailed)
		return order.ID, ErrPaymentUnavailable
	}

	// 3. Update status based on response
	if status == "Authorized" {
		if err := s.repo.UpdateStatus(ctx, order.ID, model.OrderStatusPaid); err != nil {
			return order.ID, err
		}
	} else {
		// "Declined" or other failure
		if err := s.repo.UpdateStatus(ctx, order.ID, model.OrderStatusFailed); err != nil {
			return order.ID, err
		}
	}

	return order.ID, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) CancelOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Invariants check
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