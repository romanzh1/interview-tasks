package domain

import "time"

type Order struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
