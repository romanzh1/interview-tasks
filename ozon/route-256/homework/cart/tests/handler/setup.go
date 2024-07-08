package handler

import (
	"net/http"
	"net/http/httptest"

	"route256/cart/cmd/config"
	"route256/cart/internal/handler"
	"route256/cart/internal/repository"
	"route256/cart/internal/service"
	"route256/cart/internal/service/mocks"
	"route256/cart/pkg/product"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
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

	productClient, err := product.NewClient(cfg.ProductClient)
	if err != nil {
		suite.T().Fatalf("Failed to parse product client: %v", err)
	}

	lomsClientMock := mocks.NewLomsClientMock(suite.T())

	repo := repository.NewRepository()
	svc := service.NewService(repo, productClient, lomsClientMock)
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
	goleak.VerifyNone(suite.T())
}
