package transport

import (
	"net/http"

	"github.com/Oralkhan-coder/payment-service/internal/http/dto"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	srv PaymentSrv
}

func NewPaymentHandler(srv PaymentSrv) *PaymentHandler {
	return &PaymentHandler{srv: srv}
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var req dto.PaymentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.srv.ProcessPayment(c.Request.Context(), req.OrderID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.PaymentCreateResponse{
		Status:        string(payment.Status),
		TransactionID: payment.TransactionID,
	})
}

func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	orderID := c.Param("order_id")

	payment, err := h.srv.GetPaymentByOrderID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}
