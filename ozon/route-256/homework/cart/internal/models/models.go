package models

import (
	"route256/cart/pkg/loms"
)

type CartItem struct {
	SkuID int64
	Name  string
	Count uint16
	Price uint32
}

type CartRequest struct {
	UserID int64  `json:"user_id" validate:"required,gt=0"`
	SkuID  int64  `json:"sku_id" validate:"required,gt=0"`
	Count  uint16 `json:"count" validate:"required,gt=0"`
}

func CartItemToCartRequest(ci CartItem) CartRequest {
	return CartRequest{
		SkuID: ci.SkuID,
		Count: ci.Count,
	}
}

func CartRequestToCartOrder(cr []CartRequest) []loms.OrderItem {
	items := make([]loms.OrderItem, 0, len(cr))

	for _, item := range cr {
		items = append(items, loms.OrderItem{
			SkuID:    uint64(item.SkuID),
			Quantity: uint32(item.Count),
		})
	}

	return items
}
