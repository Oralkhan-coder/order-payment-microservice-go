package http

import (
	"context"

	"github.com/Oralkhan-coder/payment-service/internal/http/transport"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router         *gin.Engine
	paymentHandler *transport.PaymentHandler
}

func NewServer(srv transport.PaymentSrv) *Server {
	paymentHandler := transport.NewPaymentHandler(srv)

	router := gin.Default()

	router.POST("/payments", paymentHandler.ProcessPayment)
	router.GET("/payments/:order_id", paymentHandler.GetPaymentStatus)

	return &Server{paymentHandler: paymentHandler, router: router}
}

func (s *Server) Run(ctx context.Context, addr string) error {
	return s.router.Run(addr)
}
