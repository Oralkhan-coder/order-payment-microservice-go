package service

import (
	"context"

	"github.com/Oralkhan-coder/payment-service/internal/model"
	"github.com/google/uuid"
)

type PaymentRepo interface {
	Create(ctx context.Context, p model.Payment) error
	GetByOrderID(ctx context.Context, orderID string) (*model.Payment, error)
}

type PaymentService struct {
	repo PaymentRepo
}

func NewPaymentService(repo PaymentRepo) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, orderID string, amount int64) (model.Payment, error) {
	payment := model.Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		TransactionID: uuid.New().String(),
		Amount:        amount,
		Status:        model.PaymentStatusAuthorized,
	}
	err := s.repo.Create(ctx, payment)
	return payment, err
}

func (s *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	return s.repo.GetByOrderID(ctx, orderID)
}
