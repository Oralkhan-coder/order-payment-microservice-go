package service

import (
	"context"
	"errors"

	paymentv1 "github.com/Oralkhan-coder/order-payment-proto-generation/payment/v1"
	"github.com/Oralkhan-coder/payment-service/internal/model"
	"github.com/google/uuid"
)

type PaymentRepo interface {
	Create(ctx context.Context, p model.Payment) error
	GetByOrderID(ctx context.Context, orderID string) (*model.Payment, error)
	FindByAmountRange(ctx context.Context, min, max int64) ([]*model.Payment, error)
}

type PaymentService struct {
	repo PaymentRepo
}

func NewPaymentService(repo PaymentRepo) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, orderID string, amount int64) (model.Payment, error) {
	status := model.PaymentStatusAuthorized
	if amount > 100000 {
		status = model.PaymentStatusDeclined
	}

	payment := model.Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		TransactionID: uuid.New().String(),
		Amount:        amount,
		Status:        status,
	}

	err := s.repo.Create(ctx, payment)
	return payment, err
}

func (s *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	return s.repo.GetByOrderID(ctx, orderID)
}

func (s *PaymentService) ListPayments(ctx context.Context, req *paymentv1.ListPaymentsRequest) ([]*model.Payment, error) {
	if req.MinAmount > 0 && req.MaxAmount > 0 && req.MinAmount > req.MaxAmount {
		return nil, errors.New("min_amount cannot be greater than max_amount")
	}

	return s.repo.FindByAmountRange(ctx, req.MinAmount, req.MaxAmount)
}
