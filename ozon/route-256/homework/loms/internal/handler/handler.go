package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"route256/loms/internal/models"
	"route256/loms/internal/service"
	"route256/loms/proto"
)

var _ proto.LomsServiceServer = (*Server)(nil)

type Server struct {
	Service  *service.Service
	Validate *validator.Validate
	Tracer   trace.Tracer
	proto.UnimplementedLomsServiceServer
}

func NewLomsService(service *service.Service, tracer trace.Tracer) *Server {
	return &Server{
		Service:  service,
		Validate: validator.New(),
		Tracer:   tracer,
	}
}

func (s *Server) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	ctx, span := s.Tracer.Start(ctx, "CreateOrder")
	defer span.End()

	if err := s.Validate.Struct(req); err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	items := make([]models.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = models.OrderItem{
			SkuID:    item.SkuId,
			Quantity: item.Quantity,
		}
	}

	orderID, err := s.Service.CreateOrder(ctx, req.User, items)
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	return &proto.CreateOrderResponse{OrderID: orderID}, nil
}

func (s *Server) GetOrder(ctx context.Context, req *proto.OrderInfoRequest) (*proto.OrderInfoResponse, error) {
	ctx, span := s.Tracer.Start(ctx, "GetOrder")
	defer span.End()

	if err := s.Validate.Struct(req); err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	order, err := s.Service.GetOrder(ctx, req.OrderID)
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	items := make([]*proto.OrderItem, len(order.Items))
	for i, item := range order.Items {
		items[i] = item.ToProto()
	}

	return &proto.OrderInfoResponse{
		Status: string(order.Status),
		User:   order.UserID,
		Items:  items,
	}, nil
}

func (s *Server) GetStockInfo(ctx context.Context, req *proto.StockInfoRequest) (*proto.StockInfoResponse, error) {
	ctx, span := s.Tracer.Start(ctx, "GetStockInfo")
	defer span.End()

	if err := s.Validate.Struct(req); err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	count, err := s.Service.GetStockInfo(ctx, req.Sku)
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.NotFound, "stock not found: %v", err)
	}

	return &proto.StockInfoResponse{Count: count}, nil
}

func (s *Server) PayOrder(ctx context.Context, req *proto.PayOrderRequest) (*emptypb.Empty, error) {
	ctx, span := s.Tracer.Start(ctx, "OrderPay")
	defer span.End()

	if err := s.Validate.Struct(req); err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	err := s.Service.OrderPay(ctx, req.OrderID)
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) CancelOrder(ctx context.Context, req *proto.OrderInfoRequest) (*emptypb.Empty, error) {
	ctx, span := s.Tracer.Start(ctx, "CancelOrder")
	defer span.End()

	if err := s.Validate.Struct(req); err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	err := s.Service.CancelOrder(ctx, req.OrderID)
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	return &emptypb.Empty{}, nil
}
