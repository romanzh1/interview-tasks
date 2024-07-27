package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/validator.v2"

	"week-2-workshop/cart/internal/app/defenitions"
	"week-2-workshop/cart/internal/domain"
	"week-2-workshop/cart/internal/service/cart/item/add"
)

type (
	addItemCommand interface {
		AddItem(ctx context.Context, userID int64, item domain.Item) error
	}

	AddHandler struct {
		name           string
		addItemCommand addItemCommand
	}
	addRequest struct {
		// request body
		Count uint16 `son:"count" validate:"nonzero"`

		// url params
		SKU  int64 `validate:"nonzero"`
		User int64 `validate:"nonzero"`
	}
)

func NewAddHandler(command addItemCommand, name string) *AddHandler {
	return &AddHandler{
		name:           name,
		addItemCommand: command,
	}
}

func (h *AddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		request *addRequest
		err     error
	)

	if request, err = h.gerRequestData(r); err != nil {
		GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}

	if err = validator.Validate(request); err != nil {
		GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}

	err = h.addItemCommand.AddItem(
		ctx,
		request.User,
		domain.Item{
			SKU:   uint32(request.SKU),
			Count: request.Count,
		},
	)

	if err != nil {
		if errors.Is(err, add.ErrInvalidSKU) {
			GetErrorResponse(w, h.name, fmt.Errorf("command handler failed: %w", err), http.StatusNotFound)
			return
		}
		GetErrorResponse(w, h.name, fmt.Errorf("command handler failed: %w", err), http.StatusInternalServerError)
		return
	}

	GetSuccessResponse(w)
}

func (h *AddHandler) gerRequestData(r *http.Request) (request *addRequest, err error) {
	request = &addRequest{}
	if err = json.NewDecoder(r.Body).Decode(request); err != nil {
		return
	}

	if request.User, err = strconv.ParseInt(r.PathValue(defenitions.ParamUserID), 10, 64); err != nil {
		return
	}

	if request.SKU, err = strconv.ParseInt(r.PathValue(defenitions.ParamSkuID), 10, 64); err != nil {
		return
	}

	return
}
