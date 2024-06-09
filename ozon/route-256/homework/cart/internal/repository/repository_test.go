package repository

import (
	"context"
	"testing"

	"route256/cart/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestAddItemToUserCart(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()

	err := repo.AddItemToUserCart(ctx, models.CartRequest{UserID: 1, SkuID: 1001, Count: 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, int(repo.carts[1][1001]))

	err = repo.AddItemToUserCart(ctx, models.CartRequest{UserID: 2, SkuID: 1002, Count: 30})
	assert.NoError(t, err)
	assert.Equal(t, 30, int(repo.carts[2][1002]))
}

func TestDeleteItemFromUserCart(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()

	err := repo.AddItemToUserCart(ctx, models.CartRequest{UserID: 1, SkuID: 1001, Count: 2})
	assert.NoError(t, err)

	err = repo.DeleteItemFromUserCart(ctx, 1, 1001)
	assert.NoError(t, err)
	_, err = repo.ListUserCart(ctx, 1)
	if assert.Error(t, err) {
		assert.Equal(t, models.ErrCartIsEmpty, err)
	}
}

func TestClearUserCart(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()

	err := repo.AddItemToUserCart(ctx, models.CartRequest{UserID: 1, SkuID: 1001, Count: 2})
	assert.NoError(t, err)
	err = repo.AddItemToUserCart(ctx, models.CartRequest{UserID: 1, SkuID: 1002, Count: 2})
	assert.NoError(t, err)
	err = repo.AddItemToUserCart(ctx, models.CartRequest{UserID: 2, SkuID: 1003, Count: 2})
	assert.NoError(t, err)

	err = repo.ClearUserCart(ctx, 1)
	assert.NoError(t, err)

	_, err = repo.ListUserCart(ctx, 1)
	if assert.Error(t, err) {
		assert.Equal(t, models.ErrCartIsEmpty, err)
	}
	_, err = repo.ListUserCart(ctx, 2)
	assert.NoError(t, err)
}

func TestListUserCart(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()

	err := repo.AddItemToUserCart(ctx, models.CartRequest{UserID: 1, SkuID: 1001, Count: 4})
	assert.NoError(t, err)
	err = repo.AddItemToUserCart(ctx, models.CartRequest{UserID: 1, SkuID: 1002, Count: 5})
	assert.NoError(t, err)
	err = repo.AddItemToUserCart(ctx, models.CartRequest{UserID: 2, SkuID: 1002, Count: 2})
	assert.NoError(t, err)

	items, err := repo.ListUserCart(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, items, 2)

	_, err = repo.ListUserCart(ctx, 3)
	if assert.Error(t, err) {
		assert.Equal(t, models.ErrCartIsEmpty, err)
	}
}
