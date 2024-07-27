package order

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/domain"
)

type createRequest struct {
	UserID      int64  `json:"user_id"`
	Description string `json:"description"`
}

type createResponse struct {
	OrderID int64  `json:"order_id"`
	Error   string `json:"error_message"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		req = &createRequest{}
	)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		data, _ := json.Marshal(createResponse{
			Error: err.Error(),
		})
		writeResponse(w, data, http.StatusBadRequest)

		return
	}

	// validate request todo

	orderID, err := h.orderRepo.Create(ctx, domain.Order{
		UserID:      req.UserID,
		Description: req.Description,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		data, _ := json.Marshal(createResponse{
			Error: err.Error(),
		})
		writeResponse(w, data, http.StatusInternalServerError)

		return
	}

	data, _ := json.Marshal(createResponse{
		OrderID: orderID,
	})
	writeResponse(w, data, http.StatusCreated)

	fmt.Printf("order created: %d\n", orderID)
}
