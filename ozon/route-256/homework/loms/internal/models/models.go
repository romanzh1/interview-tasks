package models

import (
	"route256/loms/proto"
)

type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "new"
	OrderStatusAwaitingPayment OrderStatus = "awaiting_payment"
	OrderStatusFailed          OrderStatus = "failed"
	OrderStatusPayed           OrderStatus = "payed"
	OrderStatusCanceled        OrderStatus = "canceled"
)

type OrderItem struct {
	SkuID    uint64
	Quantity uint32
}

func (i *OrderItem) ToProto() *proto.OrderItem {
	return &proto.OrderItem{
		SkuId:    i.SkuID,
		Quantity: i.Quantity,
	}
}

type Order struct {
	ID     int64
	UserID int64
	Items  []OrderItem
	Status OrderStatus
}

type Stock struct {
	SKU        uint32 `json:"sku"`
	TotalCount uint64 `json:"total_count"`
	Reserved   uint64 `json:"reserved"`
}
