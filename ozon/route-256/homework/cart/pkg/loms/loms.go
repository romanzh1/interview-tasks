package loms

import (
	"context"
	"fmt"

	"route256/cart/proto"
)

type OrderItem struct {
	SkuID    uint64
	Quantity uint32
}

func orderItemsToProto(items []OrderItem) []*proto.OrderItem {
	res := make([]*proto.OrderItem, len(items))

	for i, item := range items {
		res[i] = &proto.OrderItem{
			SkuId:    item.SkuID,
			Quantity: item.Quantity,
		}
	}

	return res
}

func (c *Client) CreateOrder(ctx context.Context, userID int64, items []OrderItem) (int64, error) {
	req := &proto.CreateOrderRequest{
		User:  userID,
		Items: orderItemsToProto(items),
	}

	res, err := c.client.CreateOrder(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	return res.OrderID, nil
}

func (c *Client) GetOrder(ctx context.Context, orderID int64) (*proto.OrderInfoResponse, error) {
	req := &proto.OrderInfoRequest{
		OrderID: orderID,
	}

	res, err := c.client.GetOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return res, nil
}

func (c *Client) GetStockInfo(ctx context.Context, sku uint32) (uint64, error) {
	req := &proto.StockInfoRequest{
		Sku: sku,
	}

	res, err := c.client.GetStockInfo(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("failed to get stock info: %w", err)
	}

	return res.Count, nil
}
