package order

import (
	"time"
)

type EventType string

const (
	EventOrderCreated  EventType = "order-created"
	EventOrderCanceled EventType = "order-canceled"
)

type Event struct {
	ID              int64     `json:"id"`
	EventType       EventType `json:"event"`
	OperationMoment time.Time `json:"moment"`
	//IdempotentKey   string    `json:"idempotent_key"`
}
