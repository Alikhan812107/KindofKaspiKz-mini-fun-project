package grpc

import (
	"time"

	"purchase-service/internal/repository"
	pb "purchase-service/pkg/pb/order"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	repo repository.PurchaseRepository
}

func NewOrderServer(repo repository.PurchaseRepository) *OrderServer {
	return &OrderServer{repo: repo}
}

func (s *OrderServer) SubscribeToOrderUpdates(req *pb.OrderRequest, stream pb.OrderService_SubscribeToOrderUpdatesServer) error {
	if req.OrderId == "" {
		return status.Error(codes.InvalidArgument, "order_id is required")
	}

	lastStatus := ""

	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
		}

		order, err := s.repo.GetByID(req.OrderId)
		if err != nil {
			return status.Error(codes.NotFound, "order not found")
		}

		if order.Status != lastStatus {
			lastStatus = order.Status
			err := stream.Send(&pb.OrderStatusUpdate{
				OrderId:   order.ID,
				Status:    order.Status,
				UpdatedAt: timestamppb.New(time.Now()),
			})
			if err != nil {
				return err
			}

			if order.Status == "Paid" || order.Status == "Failed" || order.Status == "Cancelled" {
				return nil
			}
		}

		time.Sleep(1 * time.Second)
	}
}
