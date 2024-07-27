package order

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/domain"
)

type listByUserIDRequest struct {
	UserID int64 `json:"user_id"`
}

type listByUserIDResponse struct {
	Orders []domain.Order `json:"orders"`
	Error  string
}

func (h *Handler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		req = &listByUserIDRequest{}
	)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		data, _ := json.Marshal(listByUserIDResponse{
			Error: err.Error(),
		})
		writeResponse(w, data, http.StatusBadRequest)

		return
	}

	// validate request todo

	orders, err := h.orderRepo.ListByUserID(ctx, req.UserID)
	if err != nil {
		data, _ := json.Marshal(listByUserIDResponse{
			Error: err.Error(),
		})
		writeResponse(w, data, http.StatusInternalServerError)

		return
	}

	data, _ := json.Marshal(listByUserIDResponse{
		Orders: orders,
	})
	writeResponse(w, data, http.StatusOK)

	fmt.Printf("order list orders by user id: %s\n", data)
}
