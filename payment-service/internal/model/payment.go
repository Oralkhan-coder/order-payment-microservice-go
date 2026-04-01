package model

type PaymentStatus string

const (
	PaymentStatusAuthorized PaymentStatus = "Authorized"
	PaymentStatusDeclined   PaymentStatus = "Declined"
)

type Payment struct {
	ID            string
	OrderID       string
	TransactionID string
	Amount        int64
	Status        PaymentStatus
}
