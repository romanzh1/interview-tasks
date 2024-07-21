package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"route256/cart/cmd/config"
	"route256/cart/internal/handler"
	"route256/cart/internal/repository"
	"route256/cart/internal/service"
	"route256/cart/pkg/loms"
	"route256/cart/pkg/product"
	"route256/libs/metrics"
	"route256/libs/tracing"
)

const serviceName = "cart"

func main() {
	slogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(slogHandler)
	slog.SetDefault(logger)

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Failed to parse config", "error", err)
		return
	}

	tp, err := tracing.InitTracer(
		cfg.Observability.TracerHost+":"+cfg.Observability.TracerPort, serviceName, tracing.GRPCTransport)
	if err != nil {
		slog.Error("Failed to initialize tracer", "error", err)
		return
	}
	defer tracing.ShutdownTracer(tp)

	tracer := tp.Tracer(serviceName)

	productClient, err := product.NewClient(cfg.ProductClient)
	if err != nil {
		slog.Error("Failed to parse product client timeout", "error", err)
		return
	}

	lomsClient := loms.NewClient(cfg.LomsClient.Host)
	err = lomsClient.Run()
	if err != nil {
		slog.Error("Failed to start loms client", "error", err)
		return
	}

	repo := repository.NewRepository()
	serv := service.NewService(repo, productClient, lomsClient)
	hand := handler.NewHandler(serv, tracer)
	router := http.NewServeMux()
	hand.RegisterRoutes(router)

	server := &http.Server{
		Addr:         ":" + cfg.ServerConfig.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		slog.Info("Starting server", slog.String("port", cfg.ServerConfig.Port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start server", "error", err)
		}
	}()

	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", metrics.MetricsHandler())

	go func() {
		slog.Info("Metrics server listening", slog.String("address", ":"+cfg.Observability.MetricPort))
		if err := http.ListenAndServe(":"+cfg.Observability.MetricPort, metricsMux); err != nil {
			slog.Error("Failed to serve metrics", "error", err)
		}
	}()

	pprofMux := http.NewServeMux()
	pprofMux.Handle("/debug/pprof/", http.HandlerFunc(http.DefaultServeMux.ServeHTTP))
	pprofServer := &http.Server{
		Addr:    ":" + cfg.Observability.PprofPort,
		Handler: pprofMux,
	}

	go func() {
		slog.Info("Pprof server listening", slog.String("address", ":"+cfg.Observability.PprofPort))
		if err := pprofServer.ListenAndServe(); err != nil {
			slog.Error("Failed to start pprof server", "error", err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exiting")
}
