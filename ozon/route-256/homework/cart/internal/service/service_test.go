package service

import (
	"context"
	"errors"
	"testing"

	"route256/cart/internal/models"
	"route256/cart/internal/service/mocks"
	"route256/cart/pkg/product"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func setupMocks(mc *minimock.Controller) (*mocks.CartRepositoryMock, *mocks.ProductClientMock, *mocks.LomsClientMock, *Service) {
	repoMock := mocks.NewCartRepositoryMock(mc)
	productClientMock := mocks.NewProductClientMock(mc)
	lomsClientMock := mocks.NewLomsClientMock(mc)
	s := NewService(repoMock, productClientMock, lomsClientMock)

	return repoMock, productClientMock, lomsClientMock, s
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestService_AddItemToUserCart(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := []struct {
		name             string
		cart             models.CartRequest
		getProductReturn *product.Product
		getProductErr    error
		expectedErr      string
	}{
		{
			name:             "Successful addition",
			cart:             models.CartRequest{UserID: 1, SkuID: 1001, Count: 2},
			getProductReturn: &product.Product{Name: "Product 1", Price: 100},
			getProductErr:    nil,
			expectedErr:      "",
		},
		{
			name:             "Product not found",
			cart:             models.CartRequest{UserID: 1, SkuID: 1002, Count: 1},
			getProductReturn: nil,
			getProductErr:    nil,
			expectedErr:      "product not found",
		},
		{
			name:             "Product client error",
			cart:             models.CartRequest{UserID: 1, SkuID: 1003, Count: 1},
			getProductReturn: nil,
			getProductErr:    errors.New("product client error"),
			expectedErr:      "failed to get product: product client error",
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock, productClientMock, lomsClientMock, s := setupMocks(mc)
			t.Cleanup(func() {
				repoMock.MinimockFinish()
				productClientMock.MinimockFinish()
				lomsClientMock.MinimockFinish()
			})

			if tt.name == "Successful addition" {
				lomsClientMock.GetStockInfoMock.Set(func(ctx context.Context, sku uint32) (uint64, error) {
					return 10, nil
				})
				repoMock.AddItemToUserCartMock.Return(nil)
			}

			productClientMock.GetProductMock.Return(tt.getProductReturn, tt.getProductErr)

			err := s.AddItemToUserCart(ctx, tt.cart)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err.Error())
			}

			if tt.expectedErr == "" {
				assert.Equal(t, uint64(1), lomsClientMock.GetStockInfoAfterCounter())
			}
		})
	}
}

func TestService_DeleteItemFromUserCart(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	repoMock, _, _, s := setupMocks(mc)
	t.Cleanup(repoMock.MinimockFinish)

	ctx := context.Background()

	repoMock.DeleteItemFromUserCartMock.Return(nil)

	err := s.DeleteItemFromUserCart(ctx, 1, 1001)
	assert.NoError(t, err)
}

func TestService_ClearUserCart(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	repoMock, _, _, s := setupMocks(mc)

	ctx := context.Background()

	repoMock.ClearUserCartMock.Return(nil)

	err := s.ClearUserCart(ctx, 1)
	assert.NoError(t, err)
}

func TestService_ListUserCart(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)

	tests := []struct {
		name               string
		userID             int64
		listUserCartReturn []models.CartItem
		listUserCartErr    error
		getProductReturns  map[uint32]*product.Product
		getProductErr      error
		expectedItems      []models.CartItem
		expectedTotalPrice uint32
		expectedErr        string
	}{
		{
			name:   "Successful listing",
			userID: 1,
			listUserCartReturn: []models.CartItem{
				{SkuID: 1001, Count: 2},
				{SkuID: 1002, Count: 1},
			},
			listUserCartErr: nil,
			getProductReturns: map[uint32]*product.Product{
				1001: {Name: "Product 1", Price: 100},
				1002: {Name: "Product 2", Price: 200},
			},
			expectedItems: []models.CartItem{
				{SkuID: 1001, Name: "Product 1", Count: 2, Price: 100},
				{SkuID: 1002, Name: "Product 2", Count: 1, Price: 200},
			},
			expectedTotalPrice: 400,
			expectedErr:        "",
		},
		{
			name:               "Repository error",
			userID:             1,
			listUserCartReturn: nil,
			listUserCartErr:    errors.New("repository error"),
			getProductReturns:  nil,
			expectedItems:      nil,
			expectedTotalPrice: 0,
			expectedErr:        "failed to list cart: repository error",
		},
		{
			name:   "Product client error",
			userID: 1,
			listUserCartReturn: []models.CartItem{
				{SkuID: 1001, Count: 2},
			},
			listUserCartErr:    nil,
			getProductErr:      errors.New("product client error"),
			expectedItems:      nil,
			expectedTotalPrice: 0,
			expectedErr:        "failed to get products: failed to get product: product client error",
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock, productClientMock, _, s := setupMocks(mc)

			t.Cleanup(func() {
				repoMock.MinimockFinish()
				productClientMock.MinimockFinish()
			})

			repoMock.ListUserCartMock.Return(tt.listUserCartReturn, tt.listUserCartErr)
			if tt.name != "Repository error" {
				productClientMock.GetProductMock.Set(func(ctx context.Context, skuID uint32) (p1 *product.Product, err error) {
					if tt.getProductErr != nil {
						return nil, tt.getProductErr
					}
					return tt.getProductReturns[skuID], nil
				})
			}

			items, totalPrice, err := s.ListUserCart(ctx, tt.userID)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedItems, items)
				assert.Equal(t, tt.expectedTotalPrice, totalPrice)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err.Error())
			}
		})
	}
}
