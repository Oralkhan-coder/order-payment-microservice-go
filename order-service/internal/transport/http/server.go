package http

import (
	"context"

	"github.com/Oralkhan-coder/order-service/internal/transport/http/handler"
	"github.com/Oralkhan-coder/order-service/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router       *gin.Engine
	orderHandler *handler.OrderHandler
}

func NewServer(srv handler.OrderSrv, limiter middleware.RateLimiter) *Server {
	orderHandler := handler.NewOrderHandler(srv)

	router := gin.Default()
	router.Use(middleware.RateLimit(limiter))

	router.POST("/order", orderHandler.CreateOrder)
	router.GET("/order/:id", orderHandler.GetOrder)
	router.POST("/order/:id/cancel", orderHandler.CancelOrder)

	return &Server{orderHandler: orderHandler, router: router}
}

func (s *Server) Run(ctx context.Context) {
	err := s.router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
