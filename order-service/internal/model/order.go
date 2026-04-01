package model

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "Pending"
	OrderStatusPaid      OrderStatus = "Paid"
	OrderStatusFailed    OrderStatus = "Failed"
	OrderStatusCancelled OrderStatus = "Cancelled"
)

type Order struct {
	ID         string
	CustomerID string
	ItemName   string
	Amount     int64 // Amount in cents (e.g., 1000 = $10.00)
	Status     OrderStatus // "Pending", "Paid", "Failed", "Cancelled"
	CreatedAt  time.Time
}

