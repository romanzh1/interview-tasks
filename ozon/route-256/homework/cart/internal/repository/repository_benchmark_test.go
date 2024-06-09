package repository

import (
	"context"
	"testing"

	"route256/cart/internal/models"
)

func BenchmarkAddItemToUserCart(b *testing.B) {
	repo := NewRepository()

	cart := models.CartRequest{
		UserID: int64(1),
		SkuID:  int64(1001),
		Count:  uint16(1),
	}
	ctx := context.Background()

	for n := 0; n < b.N; n++ {
		err := repo.AddItemToUserCart(ctx, cart)
		if err != nil {
			b.Fatalf("failed to add item to user cart: %v", err)
		}
	}
}

func BenchmarkListUserCart(b *testing.B) {
	repo := NewRepository()
	cart := models.CartRequest{
		UserID: int64(1),
		SkuID:  int64(1001),
		Count:  uint16(1),
	}
	ctx := context.Background()

	for i := int64(0); i < 100; i++ {
		repo.AddItemToUserCart(ctx, cart) //nolint:errcheck
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := repo.ListUserCart(ctx, cart.UserID)
		if err != nil {
			b.Fatalf("failed to list user cart: %v", err)
		}
	}
}
