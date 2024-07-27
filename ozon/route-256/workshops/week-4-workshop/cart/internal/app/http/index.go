package http

import (
	"net/http"
)

type IndexHandler struct{}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	GetSuccessResponseWithBody(w, []byte("Service 'cart' is online"))
}
