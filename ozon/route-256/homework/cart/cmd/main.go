package main

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"route256/cart/cmd/config"
	"route256/cart/internal/handler"
	"route256/cart/internal/repository"
	"route256/cart/internal/service"
	"route256/cart/pkg/product"
)

func main() {
	cfg, err := config.NewConfig()
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

	server := &http.Server{
		Addr:         ":" + cfg.ServerConfig.Port,
		Handler:      loggedMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	slog.Info("Starting server", slog.String("port", cfg.ServerConfig.Port))
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
	}
}
