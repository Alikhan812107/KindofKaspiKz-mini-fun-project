// Package payment defines the gRPC contract for PaymentService.
// This mirrors the payment.proto file in the protos repository.
package payment

import (
	"context"

	"google.golang.org/grpc"
)

// PaymentRequest is sent by the client to process a payment.
type PaymentRequest struct {
	OrderId string `json:"order_id"`
	Amount  int64  `json:"amount"`
}

// PaymentResponse is returned after processing a payment.
type PaymentResponse struct {
	PaymentId     string `json:"payment_id"`
	TransactionId string `json:"transaction_id"`
	Status        string `json:"status"`
}

const PaymentService_ProcessPayment_FullMethodName = "/payment.PaymentService/ProcessPayment"

// PaymentServiceClient is the client API for PaymentService.
type PaymentServiceClient interface {
	ProcessPayment(ctx context.Context, in *PaymentRequest, opts ...grpc.CallOption) (*PaymentResponse, error)
}

type paymentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentServiceClient(cc grpc.ClientConnInterface) PaymentServiceClient {
	return &paymentServiceClient{cc}
}

func (c *paymentServiceClient) ProcessPayment(ctx context.Context, in *PaymentRequest, opts ...grpc.CallOption) (*PaymentResponse, error) {
	out := new(PaymentResponse)
	err := c.cc.Invoke(ctx, PaymentService_ProcessPayment_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentServiceServer is the server API for PaymentService.
type PaymentServiceServer interface {
	ProcessPayment(context.Context, *PaymentRequest) (*PaymentResponse, error)
	mustEmbedUnimplementedPaymentServiceServer()
}

// UnimplementedPaymentServiceServer must be embedded to have forward-compatible implementations.
type UnimplementedPaymentServiceServer struct{}

func (UnimplementedPaymentServiceServer) ProcessPayment(context.Context, *PaymentRequest) (*PaymentResponse, error) {
	return nil, nil
}
func (UnimplementedPaymentServiceServer) mustEmbedUnimplementedPaymentServiceServer() {}

func RegisterPaymentServiceServer(s grpc.ServiceRegistrar, srv PaymentServiceServer) {
	s.RegisterService(&PaymentService_ServiceDesc, srv)
}

var PaymentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "payment.PaymentService",
	HandlerType: (*PaymentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProcessPayment",
			Handler:    _PaymentService_ProcessPayment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "payment/payment.proto",
}

func _PaymentService_ProcessPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).ProcessPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_ProcessPayment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).ProcessPayment(ctx, req.(*PaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}
