package grpc

import (
	"time"

	"purchase-service/internal/repository"
	pb "purchase-service/pkg/pb/order"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	repo repository.PurchaseRepository
}

func NewOrderServer(repo repository.PurchaseRepository) *OrderServer {
	return &OrderServer{repo: repo}
}

// SubscribeToOrderUpdates streams order status updates to the client.
// It polls the database every second and pushes a message whenever the status changes.
func (s *OrderServer) SubscribeToOrderUpdates(req *pb.OrderRequest, stream pb.OrderServiceServer_SubscribeToOrderUpdatesServer) error {
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
				UpdatedAt: time.Now().Format(time.RFC3339),
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
