package http

import (
	"order-service/internal/http/transport"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	orderHandler *transport.OrderHandler
}

func NewServer(order Order) *Server {
	orderHandler := transport.NewOrderHandler(order)

	router := gin.Default()

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