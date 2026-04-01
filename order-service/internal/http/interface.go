package http

import "order-service/internal/http/transport"

type Order interface {
	transport.OrderSrv
}