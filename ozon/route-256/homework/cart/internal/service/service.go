package service

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"route256/cart/internal/models"
	"route256/cart/pkg/errgroup"
	"route256/cart/pkg/loms"
	"route256/cart/pkg/product"
	"route256/libs/metrics"
)

const countGoRoutines = 10

//go:generate minimock -i cartRepository -o ./mocks/ -s _mock.go
//go:generate minimock -i productClient -o ./mocks/ -s _mock.go
//go:generate minimock -i lomsClient -o ./mocks/ -s _mock.go
type cartRepository interface {
	AddItemToUserCart(ctx context.Context, cart models.CartRequest) error
	DeleteItemFromUserCart(ctx context.Context, userID, skuID int64) error
	ClearUserCart(ctx context.Context, userID int64) error
	ListUserCart(ctx context.Context, userID int64) ([]models.CartItem, error)
}

type productClient interface {
	GetProduct(ctx context.Context, skuID uint32) (*product.Product, error)
}

type lomsClient interface {
	CreateOrder(ctx context.Context, userID int64, items []loms.OrderItem) (int64, error)
	GetStockInfo(ctx context.Context, sku uint32) (uint64, error)
}

type Service struct {
	repo          cartRepository
	productClient productClient
	lomsClient    lomsClient
	tracer        trace.Tracer
}

func NewService(repo cartRepository, productClient productClient, lomsClient lomsClient) *Service {
	tracer := otel.Tracer("cartService")
	return &Service{repo: repo, productClient: productClient, lomsClient: lomsClient, tracer: tracer}
}

func (s *Service) CreateOrder(ctx context.Context, userID int64) (int64, error) {
	ctx, span := s.tracer.Start(ctx, "Service.CreateOrder")
	defer span.End()

	items, err := s.repo.ListUserCart(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to list user cart: %w", err)
	}

	stockInfo := make([]models.CartRequest, len(items))
	for i, item := range items {
		metrics.ExternalRequests.WithLabelValues("GetStockInfo", "count").Inc()
		count, err := s.lomsClient.GetStockInfo(ctx, uint32(item.SkuID))
		if err != nil {
			return 0, fmt.Errorf("failed to get stock info: %w", err)
		}

		stockInfo[i].Count = uint16(count)
		if stockInfo[i].Count < item.Count || stockInfo[i].Count == 0 {
			return 0, fmt.Errorf("stock is not enough")
		}

		stockInfo[i].SkuID = item.SkuID
		stockInfo[i].UserID = userID
		stockInfo[i].Count = item.Count
	}

	metrics.ExternalRequests.WithLabelValues("CreateOrder", "count").Inc()
	orderID, err := s.lomsClient.CreateOrder(ctx, userID, models.CartRequestToCartOrder(stockInfo))
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	err = s.repo.ClearUserCart(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to clear user cart: %w", err)
	}

	return orderID, nil
}

func (s *Service) AddItemToUserCart(ctx context.Context, cart models.CartRequest) error {
	ctx, span := s.tracer.Start(ctx, "Service.AddItemToUserCart")
	defer span.End()

	metrics.ExternalRequests.WithLabelValues("GetProduct", "count").Inc()
	checkProduct, err := s.productClient.GetProduct(ctx, uint32(cart.SkuID))
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	if checkProduct == nil {
		return fmt.Errorf("product not found")
	}

	metrics.ExternalRequests.WithLabelValues("GetStockInfo", "count").Inc()
	count, err := s.lomsClient.GetStockInfo(ctx, uint32(cart.SkuID))
	if err != nil {
		return fmt.Errorf("failed to get stock info: %w", err)
	}

	if uint16(count) < cart.Count || count == 0 {
		return fmt.Errorf("stock is not enough")
	}

	return s.repo.AddItemToUserCart(ctx, cart)
}

func (s *Service) DeleteItemFromUserCart(ctx context.Context, userID, skuID int64) error {
	return s.repo.DeleteItemFromUserCart(ctx, userID, skuID)
}

func (s *Service) ClearUserCart(ctx context.Context, userID int64) error {
	return s.repo.ClearUserCart(ctx, userID)
}

func (s *Service) ListUserCart(ctx context.Context, userID int64) ([]models.CartItem, uint32, error) {
	ctx, span := s.tracer.Start(ctx, "Service.ListUserCart")
	defer span.End()

	items, err := s.repo.ListUserCart(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list cart: %w", err)
	}

	var (
		cartItems  = make([]models.CartItem, len(items))
		totalPrice uint32
		mu         sync.Mutex
		g          = errgroup.NewGroup(countGoRoutines)
	)

	for i, item := range items {
		g.Go(func() error {
			metrics.ExternalRequests.WithLabelValues("GetProduct", "count").Inc()
			p, err := s.productClient.GetProduct(ctx, uint32(item.SkuID))
			if err != nil {
				return fmt.Errorf("failed to get product: %w", err)
			}
			if p == nil {
				return nil
			}

			mu.Lock()
			cartItems[i] = models.CartItem{
				SkuID: item.SkuID,
				Name:  p.Name,
				Count: item.Count,
				Price: p.Price,
			}
			totalPrice += p.Price * uint32(item.Count)
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, 0, fmt.Errorf("failed to get products: %w", err)
	}

	sort.Slice(cartItems, func(i, j int) bool {
		return cartItems[i].SkuID < cartItems[j].SkuID
	})

	return cartItems, totalPrice, nil
}
