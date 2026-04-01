package repository

import (
	"context"
	"time"

	"github.com/Oralkhan-coder/order-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) Create(ctx context.Context, o model.Order) error {
	q := `
		INSERT INTO orders (id, customer_id, item_name, amount, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, q, o.ID, o.CustomerID, o.ItemName, o.Amount, o.Status, time.Now())
	return err
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	q := `
		SELECT id, customer_id, item_name, amount, status, created_at
		FROM orders
		WHERE id = $1
	`
	var o model.Order
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&o.ID, &o.CustomerID, &o.ItemName, &o.Amount, &o.Status, &o.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error {
	q := `
		UPDATE orders
		SET status=$1	
		WHERE id=$2
	`
	_, err := r.pool.Exec(ctx, q, status, id)
	return err
}
