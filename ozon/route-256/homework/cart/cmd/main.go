package main

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	handler "route256/cart/internal/handler"
	"route256/cart/internal/repository"
	"route256/cart/internal/service"
	"route256/cart/pkg/product"
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatalf("failed to parse config: %s", err)
	}

	productClientTimeout, err := time.ParseDuration(cfg.ProductClient.Timeout + "s")
	if err != nil {
		log.Fatalf("failed to parse product client timeout: %s", err)
	}

	router := http.NewServeMux()
	productClient := product.NewClient(cfg.ProductClient.Host, cfg.ProductClient.Token, productClientTimeout)

	repo := repository.NewRepository()
	serv := service.NewService(repo, productClient)
	hand := handler.NewHandler(serv)
	hand.RegisterRoutes(router)
	loggedMux := handler.LoggingMiddleware(router)

	slog.Info("Starting server", slog.String("port", cfg.ServerConfig.Port))
	if err := http.ListenAndServe(":"+cfg.ServerConfig.Port, loggedMux); err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
	}
}
