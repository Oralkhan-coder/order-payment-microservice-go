package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) Create(ctx context.Context, o model.Order) (int, error) {
	q := `
		INSERT INTO orders (customer_id, item_name, amount)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id int
	err := r.pool.QueryRow(ctx, q, o.CustomerID, o.ItemName, o.Amount).Scan(&id)
	return id, err
}

func (r *OrderRepository) GetByID(ctx context.Context, id int) (*model.Order, error) {
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

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int, status model.OrderStatus) error {
	q := `
		UPDATE orders
		SET status=$1	
		WHERE id=$2
	`
	_, err := r.pool.Exec(ctx, q, status, id)
	return err
}
