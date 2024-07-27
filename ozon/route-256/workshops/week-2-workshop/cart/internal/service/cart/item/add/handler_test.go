package add

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
	"week-2-workshop/cart/internal/domain"
	"week-2-workshop/cart/internal/service/cart/item/add/mock"
)

func TestAddItem(t *testing.T) {
	ctx := context.Background()
	var (
		userID int64
		item   domain.Item
	)
	userID = 100
	item = domain.Item{
		SKU:   1000,
		Count: 10,
	}

	ctrl := minimock.NewController(t)
	productMock := mock.NewProductServiceMock(ctrl)
	repMock := mock.NewRepositoryMock(ctrl)
	addHandler := New(repMock, productMock)

	product := domain.Product{
		Name:  "Книга",
		Price: 200,
	}

	productMock.GetProductInfoMock.Expect(ctx, 1000).Return(&product, nil)
	repMock.AddItemMock.Expect(ctx, 100, item).Return(nil)

	err := addHandler.AddItem(ctx, userID, item)
	require.NoError(t, err)

	t.Log(productMock.GetProductInfoMock.Expect(ctx, 1000).Return(&product, nil).GetProductInfoBeforeCounter())
}

func TestAddItemTable(t *testing.T) {
	ctx := context.Background()

	type data struct {
		name    string
		userID  int64
		item    domain.Item
		product *domain.Product
		wantErr error
	}

	testData := []data{{
		name:   "valid add item",
		userID: 123,
		item: domain.Item{
			SKU:   100,
			Count: 2,
		},
		product: &domain.Product{
			Name:  "Книга",
			Price: 400,
		},
		wantErr: nil,
	},
		{
			name:   "product not found",
			userID: 123,
			item: domain.Item{
				SKU:   100,
				Count: 2,
			},
			product: nil,
			wantErr: ErrInvalidSKU,
		}}

	ctrl := minimock.NewController(t)
	productMock := mock.NewProductServiceMock(ctrl)
	repMock := mock.NewRepositoryMock(ctrl)
	addHandler := New(repMock, productMock)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			productMock.GetProductInfoMock.Expect(ctx, tt.item.SKU).Return(tt.product, nil)
			repMock.AddItemMock.Expect(ctx, tt.userID, tt.item).Return(nil)

			err := addHandler.AddItem(ctx, tt.userID, tt.item)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestAddItemTableWithPrepare(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		productMock *mock.ProductServiceMock
		repMock     *mock.RepositoryMock
	}
	type data struct {
		name    string
		userID  int64
		item    domain.Item
		prepare func(f *fields)
		wantErr error
	}

	testData := []data{{
		name:   "valid add item",
		userID: 123,
		item: domain.Item{
			SKU:   100,
			Count: 2,
		},
		prepare: func(f *fields) {
			f.productMock.GetProductInfoMock.Expect(ctx, 100).Return(&domain.Product{
				Name:  "Книга",
				Price: 300,
			}, nil)
			f.repMock.AddItemMock.Expect(ctx, 123, domain.Item{
				SKU:   100,
				Count: 2,
			}).Return(nil)
		},
		wantErr: nil,
	},
		{
			name:   "product not found",
			userID: 123,
			item: domain.Item{
				SKU:   100,
				Count: 2,
			},
			prepare: func(f *fields) {
				f.productMock.GetProductInfoMock.Expect(ctx, 100).Return(nil, nil)
			},
			wantErr: ErrInvalidSKU,
		}}

	ctrl := minimock.NewController(t)
	fieldsForTableTest := fields{
		productMock: mock.NewProductServiceMock(ctrl),
		repMock:     mock.NewRepositoryMock(ctrl),
	}

	addHandler := New(fieldsForTableTest.repMock, fieldsForTableTest.productMock)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(&fieldsForTableTest)
			err := addHandler.AddItem(ctx, tt.userID, tt.item)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestAddItemV2(t *testing.T) {
	ctx := minimock.AnyContext
	var (
		userID int64
	)
	userID = 100
	item1 := domain.Item{
		SKU:   1000,
		Count: 10,
	}
	item2 := domain.Item{
		SKU:   2000,
		Count: 20,
	}

	ctrl := minimock.NewController(t)
	productMock := mock.NewProductServiceMock(ctrl)
	repMock := mock.NewRepositoryMock(ctrl)
	addHandler := New(repMock, productMock)

	product := domain.Product{
		Name:  "Книга",
		Price: 200,
	}

	productMock.GetProductInfoMock.When(ctx, 1000).Then(&product, nil).
		GetProductInfoMock.When(ctx, 2000).Then(nil, nil)

	//productMock.GetProductInfoMock.Set(func(ctx context.Context, sku uint32) (pp1 *domain.Product, err error) {
	//	switch sku {
	//	case 1000:
	//		return &product, nil
	//	case 2000:
	//		return nil, nil
	//	default:
	//		return nil, nil
	//	}
	//})

	err := addHandler.AddItemV2(ctx, userID, item1, item2)
	require.ErrorIs(t, err, ErrInvalidSKU)
}
