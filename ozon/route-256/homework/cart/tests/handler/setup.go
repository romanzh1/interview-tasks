package handler

import (
	"net/http"
	"net/http/httptest"
	"time"

	"route256/cart/cmd/config"
	"route256/cart/internal/handler"
	"route256/cart/internal/repository"
	"route256/cart/internal/service"
	"route256/cart/pkg/product"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	server  *httptest.Server
	repo    *repository.Repository
	service *service.Service
	handler *handler.Handler
}

func (suite *TestSuite) SetupSuite() {
	err := godotenv.Load("../../.env")
	if err != nil {
		suite.T().Fatalf("Error loading .env file: %v", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		suite.T().Fatalf("Failed to parse config: %v", err)
	}

	productClientTimeout, err := time.ParseDuration(cfg.ProductClient.Timeout + "s")
	if err != nil {
		suite.T().Fatalf("Failed to parse product client timeout: %v", err)
	}

	productClient := product.NewClient(cfg.ProductClient.Host, cfg.ProductClient.Token, productClientTimeout)
	repo := repository.NewRepository()
	svc := service.NewService(repo, productClient)
	handlers := handler.NewHandler(svc)

	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux)

	suite.repo = repo
	suite.service = svc
	suite.handler = handlers

	server := httptest.NewServer(mux)

	suite.server = server
}

func (suite *TestSuite) TearDownSuite() {
	suite.server.Close()
}
