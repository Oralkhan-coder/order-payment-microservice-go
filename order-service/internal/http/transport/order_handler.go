package transport

import (
	"context"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	srv OrderSrv
}

func NewOrderHandler(srv OrderSrv) *OrderHandler {
	return &OrderHandler{srv: srv}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.srv.CreateOrder(c.Request.Context(), req.CustomerID, req.ItemName, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order_id": id})	
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	order, err := h.srv.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)		
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")

	err := h.srv.CancelOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})	
}