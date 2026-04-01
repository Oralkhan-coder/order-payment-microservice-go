package repository

import (
	"context"

	"github.com/Oralkhan-coder/payment-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{pool: pool}
}

func (r *PaymentRepository) Create(ctx context.Context, p model.Payment) error {
	q := `
		INSERT INTO payments (id, order_id, transaction_id, amount, status)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, q, p.ID, p.OrderID, p.TransactionID, p.Amount, p.Status)
	return err
}

func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	q := `
		SELECT id, order_id, transaction_id, amount, status
		FROM payments
		WHERE order_id = $1
	`
	var p model.Payment
	err := r.pool.QueryRow(ctx, q, orderID).Scan(
		&p.ID, &p.OrderID, &p.TransactionID, &p.Amount, &p.Status,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
