package domain

type ListItem struct {
	SKU     uint32
	Count   uint16
	Product Product
}

type Item struct {
	SKU   uint32
	Count uint16
}
