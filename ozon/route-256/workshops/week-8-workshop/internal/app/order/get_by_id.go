package order

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/domain"
)

type getByIDRequest struct {
	OrderID int64 `json:"order_id"`
}

type getByIDResponse struct {
	Order *domain.Order `json:"order"`
	Error string
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		req = &getByIDRequest{}
	)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		data, _ := json.Marshal(getByIDResponse{
			Error: err.Error(),
		})
		writeResponse(w, data, http.StatusBadRequest)

		return
	}

	// validate request todo

	order, err := h.orderRepo.GetByOrderID(ctx, req.OrderID)
	if err != nil {
		data, _ := json.Marshal(getByIDResponse{
			Error: err.Error(),
		})
		writeResponse(w, data, http.StatusInternalServerError)

		return
	}

	data, _ := json.Marshal(getByIDResponse{
		Order: order,
	})
	writeResponse(w, data, http.StatusOK)

	fmt.Printf("order get: %v+\n", order)
}
