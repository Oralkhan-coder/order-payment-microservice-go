package handler

import (
	"time"

	orderv1 "github.com/Oralkhan-coder/order-payment-proto-generation/order/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCOrderHandler struct {
	orderv1.UnimplementedOrderServiceServer
	srv OrderSrv
}

func NewGRPCOrderHandler(srv OrderSrv) *GRPCOrderHandler {
	return &GRPCOrderHandler{srv: srv}
}

func (s *GRPCOrderHandler) SubscribeToOrderUpdates(req *orderv1.OrderRequest, stream orderv1.OrderService_SubscribeToOrderUpdatesServer) error {
	if req.OrderId == "" {
		return status.Error(codes.InvalidArgument, "order_id required")
	}

	var lastStatus string
	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
		}

		currentStatus, err := s.srv.GetOrderStatus(stream.Context(), req.OrderId)

		if err != nil {
			return status.Error(codes.NotFound, "order not found")
		}

		if string(*currentStatus) != lastStatus {
			lastStatus = string(*currentStatus)
			if err := stream.Send(&orderv1.OrderStatusUpdate{
				OrderId: req.OrderId,
				Status:  string(*currentStatus),
				Message: "Status updated to " + string(*currentStatus),
			}); err != nil {
				return err
			}
		}

		time.Sleep(2 * time.Second)
	}
}
