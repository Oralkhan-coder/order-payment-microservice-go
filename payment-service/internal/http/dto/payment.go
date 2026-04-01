package dto

type PaymentCreateRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	Amount  int64  `json:"amount" binding:"required"`
}

type PaymentCreateResponse struct {
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id"`
}
