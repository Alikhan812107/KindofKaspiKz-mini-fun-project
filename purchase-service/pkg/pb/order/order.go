// Package order defines the gRPC contract for OrderService.
// This mirrors the order.proto file in the protos repository.
package order

import (
	"context"

	"google.golang.org/grpc"
)

// OrderRequest is sent by the client to subscribe to order updates.
type OrderRequest struct {
	OrderId string `json:"order_id"`
}

// OrderStatusUpdate is streamed back to the client.
type OrderStatusUpdate struct {
	OrderId   string `json:"order_id"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

const OrderService_SubscribeToOrderUpdates_FullMethodName = "/order.OrderService/SubscribeToOrderUpdates"

// OrderServiceServer_SubscribeToOrderUpdatesServer is the stream type for SubscribeToOrderUpdates.
type OrderServiceServer_SubscribeToOrderUpdatesServer = grpc.ServerStreamingServer[OrderStatusUpdate]

// OrderServiceServer is the server API for OrderService.
type OrderServiceServer interface {
	SubscribeToOrderUpdates(*OrderRequest, OrderServiceServer_SubscribeToOrderUpdatesServer) error
	mustEmbedUnimplementedOrderServiceServer()
}

// UnimplementedOrderServiceServer must be embedded to have forward-compatible implementations.
type UnimplementedOrderServiceServer struct{}

func (UnimplementedOrderServiceServer) SubscribeToOrderUpdates(*OrderRequest, OrderServiceServer_SubscribeToOrderUpdatesServer) error {
	return nil
}
func (UnimplementedOrderServiceServer) mustEmbedUnimplementedOrderServiceServer() {}

// OrderServiceClient is the client API for OrderService.
type OrderServiceClient interface {
	SubscribeToOrderUpdates(ctx context.Context, in *OrderRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[OrderStatusUpdate], error)
}

type orderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderServiceClient(cc grpc.ClientConnInterface) OrderServiceClient {
	return &orderServiceClient{cc}
}

func (c *orderServiceClient) SubscribeToOrderUpdates(ctx context.Context, in *OrderRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[OrderStatusUpdate], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &OrderService_ServiceDesc.Streams[0], OrderService_SubscribeToOrderUpdates_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[OrderRequest, OrderStatusUpdate]{ClientStream: stream}
	if err := x.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

func RegisterOrderServiceServer(s grpc.ServiceRegistrar, srv OrderServiceServer) {
	s.RegisterService(&OrderService_ServiceDesc, srv)
}

var OrderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "order.OrderService",
	HandlerType: (*OrderServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeToOrderUpdates",
			Handler:       _OrderService_SubscribeToOrderUpdates_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "order/order.proto",
}

func _OrderService_SubscribeToOrderUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(OrderRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OrderServiceServer).SubscribeToOrderUpdates(m, &grpc.GenericServerStream[OrderRequest, OrderStatusUpdate]{ServerStream: stream})
}
