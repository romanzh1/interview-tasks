package service

import (
	"errors"
	"testing"

	"route256/cart/internal/models"
	"route256/cart/internal/service/mocks"
	"route256/cart/pkg/product"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func setupMocks(mc *minimock.Controller) (*mocks.CartRepositoryMock, *mocks.ProductClientMock, *Service) {
	repoMock := mocks.NewCartRepositoryMock(mc)
	productClientMock := mocks.NewProductClientMock(mc)
	s := NewService(repoMock, productClientMock)

	return repoMock, productClientMock, s
}

func TestService_AddItemToUserCart(t *testing.T) {
	mc := minimock.NewController(t)

	repoMock, productClientMock, s := setupMocks(mc)

	tests := []struct {
		name             string
		skuID            int64
		count            uint16
		getProductReturn *product.Product
		getProductErr    error
		expectedErr      string
	}{
		{
			name:             "Successful addition",
			skuID:            1001,
			count:            2,
			getProductReturn: &product.Product{Name: "Product 1", Price: 100},
			getProductErr:    nil,
			expectedErr:      "",
		},
		{
			name:             "Product not found",
			skuID:            1002,
			count:            1,
			getProductReturn: nil,
			getProductErr:    nil,
			expectedErr:      "product not found",
		},
		{
			name:             "Product client error",
			skuID:            1003,
			count:            1,
			getProductReturn: nil,
			getProductErr:    errors.New("product client error"),
			expectedErr:      "failed to get product: product client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productClientMock.GetProductMock.Return(tt.getProductReturn, tt.getProductErr)
			repoMock.AddItemToUserCartMock.Return(nil)

			err := s.AddItemToUserCart(1, tt.skuID, tt.count)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err.Error())
			}
		})
	}
}

func TestService_DeleteItemFromUserCart(t *testing.T) {
	mc := minimock.NewController(t)

	repoMock, _, s := setupMocks(mc)

	repoMock.DeleteItemFromUserCartMock.Return(nil)

	err := s.DeleteItemFromUserCart(1, 1001)
	assert.NoError(t, err)
}

func TestService_ClearUserCart(t *testing.T) {
	mc := minimock.NewController(t)

	repoMock, _, s := setupMocks(mc)

	repoMock.ClearUserCartMock.Return(nil)

	err := s.ClearUserCart(1)
	assert.NoError(t, err)
}

func TestService_ListUserCart(t *testing.T) {
	mc := minimock.NewController(t)

	repoMock, productClientMock, s := setupMocks(mc)

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
			expectedErr:        "failed to get p: product client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock.ListUserCartMock.Return(tt.listUserCartReturn, tt.listUserCartErr)
			productClientMock.GetProductMock.Set(func(skuID uint32) (p1 *product.Product, err error) {
				if tt.getProductErr != nil {
					return nil, tt.getProductErr
				}
				return tt.getProductReturns[skuID], nil
			})

			items, totalPrice, err := s.ListUserCart(tt.userID)
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
