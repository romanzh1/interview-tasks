package item_suite

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"week-2-workshop/cart/internal/clients/product"
	"week-2-workshop/cart/internal/domain"
	"week-2-workshop/cart/internal/repository/memory_cart_repo"
	"week-2-workshop/cart/internal/service/cart/item/add"
)

type ItemS struct {
	suite.Suite
	addHandler *add.Handler
}

func (s *ItemS) SetupSuite() {
	productClient, err := product.New("http://route256.pavl.uk:8080", "testtoken")
	require.NoError(s.T(), err)

	storage := memory_cart_repo.NewMemoryStorage()

	s.addHandler = add.New(storage, productClient)
}

func (s *ItemS) TestAddItem() {
	ctx := context.Background()

	userID := int64(123)
	item := domain.Item{
		SKU:   773297411,
		Count: 2,
	}

	err := s.addHandler.AddItem(ctx, userID, item)
	//require.ErrorIs(s.T(), err, add.ErrInvalidSKU)
	require.NoError(s.T(), err)
}
