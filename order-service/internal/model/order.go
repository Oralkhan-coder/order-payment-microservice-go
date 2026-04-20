package model

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusFailed    OrderStatus = "failed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID         string
	CustomerID string
	ItemName   string
	Amount     int64       // Amount in cents (e.g., 1000 = $10.00)
	Status     OrderStatus // "Pending", "Paid", "Failed", "Cancelled"
	CreatedAt  time.Time
}
