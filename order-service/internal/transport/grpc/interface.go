package grpc

import "github.com/Oralkhan-coder/order-service/internal/transport/grpc/handler"

type Order interface {
	handler.OrderSrv
}
