package order

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/domain"
)

type listByIDRequest struct {
	OrderIDs []int64 `json:"order_ids"`
}

type listByIDResponse struct {
	Orders []domain.Order `json:"orders"`
	Error  string
}

func (h *Handler) ListByID(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		req = &listByIDRequest{}
	)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		data, _ := json.Marshal(listByIDResponse{
			Error: err.Error(),
		})
		writeResponse(w, data, http.StatusBadRequest)

		return
	}

	// validate request todo

	orders, err := h.orderRepo.ListByID(ctx, req.OrderIDs)
	if err != nil {
		data, _ := json.Marshal(listByIDResponse{
			Error: err.Error(),
		})
		writeResponse(w, data, http.StatusInternalServerError)

		return
	}

	data, _ := json.Marshal(listByIDResponse{
		Orders: orders,
	})
	writeResponse(w, data, http.StatusOK)

	fmt.Printf("order list orders by user id: %s\n", data)
}
