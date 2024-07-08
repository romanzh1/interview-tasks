package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"route256/libs/metrics"
	"route256/libs/tracing"
	"route256/loms/internal/handler"

	"route256/loms/cmd/config"
	"route256/loms/internal/repository"
	"route256/loms/internal/service"
	"route256/loms/proto"
)

const serviceName = "loms"

func main() {
	slogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(slogHandler)
	slog.SetDefault(logger)

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Failed to parse config", "error", err)
		return
	}

	conn, err := pgxpool.Connect(context.Background(), cfg.LomsDatabaseURL())
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return
	}
	defer conn.Close()

	tp, err := tracing.InitTracer(
		cfg.Observability.TracerHost+":"+cfg.Observability.TracerPort, serviceName)
	if err != nil {
		slog.Error("Failed to initialize tracer", "error", err)
		return
	}
	defer tracing.ShutdownTracer(tp)

	tracer := tp.Tracer(serviceName)

	txManager := repository.NewTxManager(conn)
	orderRepo := repository.NewOrderRepository(conn)
	stockRepo := repository.NewStockRepository(conn)

	svc := service.NewService(txManager, orderRepo, stockRepo)
	lomsHandler := handler.NewLomsService(svc, tracer)

	lis, err := net.Listen("tcp", ":"+cfg.ServerConfig.GRPCPort)
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		return
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(lomsHandler.UnaryServerInterceptor()),
	)

	gwMux := runtime.NewServeMux()

	proto.RegisterLomsServiceServer(grpcServer, lomsHandler)

	gwServer := &http.Server{
		Addr:    ":" + cfg.ServerConfig.HTTPPort,
		Handler: gwMux,
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = proto.RegisterLomsServiceHandlerFromEndpoint(context.Background(),
		gwMux, "localhost:"+cfg.ServerConfig.GRPCPort, opts)
	if err != nil {
		slog.Error("Failed to register gateway", "error", err)
		return
	}

	go func() {
		slog.Info("HTTP server listening", "address", gwServer.Addr)
		if err := gwServer.ListenAndServe(); err != nil {
			slog.Error("Failed to serve", "error", err)
		}
	}()

	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", metrics.MetricsHandler())

	go func() {
		slog.Info("Metrics server listening", "address", ":"+cfg.Observability.MetricPort)
		if err := http.ListenAndServe(":"+cfg.Observability.MetricPort, metricsMux); err != nil {
			slog.Error("Failed to serve metrics", "error", err)
		}
	}()

	slog.Info("Server listening", "address", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("Failed to serve", "error", err)
	}
}
