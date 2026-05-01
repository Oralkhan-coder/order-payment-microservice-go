package repository

import (
	"context"
	"fmt"

	infrastructure_postgres "github.com/Oralkhan-coder/payment-service/internal/infrastructure/postgres"
	"github.com/Oralkhan-coder/payment-service/internal/model"
)

type PaymentRepository struct {
	pool infrastructure_postgres.Pool
}

func NewPaymentRepository(pool infrastructure_postgres.Pool) *PaymentRepository {
	return &PaymentRepository{pool: pool}
}

func (r *PaymentRepository) Create(ctx context.Context, p model.Payment) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	q := `
		INSERT INTO payments (id, order_id, transaction_id, amount, status)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, q, p.ID, p.OrderID, p.TransactionID, p.Amount, p.Status)
	return err
}

func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

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

func (r *PaymentRepository) FindByAmountRange(ctx context.Context, min, max int64) ([]*model.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := "SELECT id, amount, status FROM payments WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if min > 0 {
		query += fmt.Sprintf(" AND amount >= $%d", argCount)
		args = append(args, min)
		argCount++
	}
	if max > 0 {
		query += fmt.Sprintf(" AND amount <= $%d", argCount)
		args = append(args, max)
		argCount++
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var payments []*model.Payment
	for rows.Next() {
		p := &model.Payment{}
		if err := rows.Scan(&p.ID, &p.Amount, &p.Status); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		payments = append(payments, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}
