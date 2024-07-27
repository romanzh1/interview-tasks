package add

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
	"week-4-workshop/cart/internal/domain"
)

func TestAddItem(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	productServiceMock := NewProductServiceMock(mc)
	repoMock := NewRepositoryMock(mc)

	addHandler := New(repoMock, productServiceMock)

	type inputData struct {
		userID int64
		item   domain.Item
	}
	tests := []struct {
		name      string
		inputData inputData
		wantErr   error
	}{
		{
			name: "product not found",
			inputData: inputData{
				userID: 1,
				item: domain.Item{
					SKU:   1,
					Count: 10,
				},
			},
			wantErr: ErrInvalidSKU,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			productServiceMock.GetProductInfoMock.Expect(ctx, tt.inputData.item.SKU).Return(nil, nil)

			err := addHandler.AddItem(ctx, tt.inputData.userID, tt.inputData.item)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
