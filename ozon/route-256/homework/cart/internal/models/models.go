package models

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
