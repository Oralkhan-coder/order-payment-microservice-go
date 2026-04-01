package http

import "github.com/Oralkhan-coder/order-service/internal/http/transport"

type Order interface {
	transport.OrderSrv
}