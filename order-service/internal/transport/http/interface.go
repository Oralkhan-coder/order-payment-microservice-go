package http

import (
	"github.com/Oralkhan-coder/order-service/internal/transport/http/handler"
)

type Order interface {
	handler.OrderSrv
}
