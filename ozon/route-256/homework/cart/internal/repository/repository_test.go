package repository

import (
	"testing"

	"route256/cart/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestAddItemToUserCart(t *testing.T) {
	repo := NewRepository()

	err := repo.AddItemToUserCart(1, 1001, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, int(repo.carts[1][1001]))

	err = repo.AddItemToUserCart(2, 1002, 30)
	assert.NoError(t, err)
	assert.Equal(t, 30, int(repo.carts[2][1002]))
}

func TestDeleteItemFromUserCart(t *testing.T) {
	repo := NewRepository()
	err := repo.AddItemToUserCart(1, 1001, 2)
	assert.NoError(t, err)

	err = repo.DeleteItemFromUserCart(1, 1001)
	assert.NoError(t, err)
	_, err = repo.ListUserCart(1)
	if assert.Error(t, err) {
		assert.Equal(t, models.CartIsEmpty, err)
	}
}

func TestClearUserCart(t *testing.T) {
	repo := NewRepository()
	err := repo.AddItemToUserCart(1, 1001, 2)
	assert.NoError(t, err)
	err = repo.AddItemToUserCart(1, 1002, 2)
	assert.NoError(t, err)
	err = repo.AddItemToUserCart(2, 1003, 2)
	assert.NoError(t, err)

	err = repo.ClearUserCart(1)
	assert.NoError(t, err)

	_, err = repo.ListUserCart(1)
	if assert.Error(t, err) {
		assert.Equal(t, models.CartIsEmpty, err)
	}
	_, err = repo.ListUserCart(2)
	assert.NoError(t, err)
}

func TestListUserCart(t *testing.T) {
	repo := NewRepository()
	err := repo.AddItemToUserCart(1, 1001, 4)
	assert.NoError(t, err)
	err = repo.AddItemToUserCart(1, 1002, 5)
	assert.NoError(t, err)
	err = repo.AddItemToUserCart(2, 1002, 2)
	assert.NoError(t, err)

	items, err := repo.ListUserCart(1)
	assert.NoError(t, err)
	assert.Len(t, items, 2)

	_, err = repo.ListUserCart(3)
	if assert.Error(t, err) {
		assert.Equal(t, models.CartIsEmpty, err)
	}
}
