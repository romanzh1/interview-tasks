package handler

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"route256/libs/metrics"
)

var grpcToHTTPStatus = map[codes.Code]string{
	codes.OK:                 "200",
	codes.Canceled:           "499",
	codes.Unknown:            "500",
	codes.InvalidArgument:    "400",
	codes.DeadlineExceeded:   "504",
	codes.NotFound:           "404",
	codes.AlreadyExists:      "409",
	codes.PermissionDenied:   "403",
	codes.ResourceExhausted:  "429",
	codes.FailedPrecondition: "412",
	codes.Aborted:            "409",
	codes.OutOfRange:         "400",
	codes.Unimplemented:      "501",
	codes.Internal:           "500",
	codes.Unavailable:        "503",
	codes.DataLoss:           "500",
	codes.Unauthenticated:    "401",
}

func (s *Server) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start).Seconds()

		var statusCode string
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				statusCode = grpcToHTTPStatus[st.Code()]
			} else {
				statusCode = "unknown"
			}
			slog.Error("gRPC method handler error", "method", info.FullMethod, "error", err)
		} else {
			statusCode = "200"
			slog.Info("gRPC method succeeded", "method", info.FullMethod)
		}

		metrics.RequestDuration.WithLabelValues(info.FullMethod, statusCode).Observe(duration)
		metrics.TotalRequests.WithLabelValues(info.FullMethod, statusCode).Inc()

		return resp, err
	}
}
