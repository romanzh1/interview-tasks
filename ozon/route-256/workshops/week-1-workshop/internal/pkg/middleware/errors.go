package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorWrapper func(writer http.ResponseWriter, request *http.Request) error

func (s ErrorWrapper) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if err := s(writer, request); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		var errorHandler struct {
			Code    int
			Message string
		}

		errorHandler.Message = err.Error()
		errorHandler.Code = 989

		buf, err := json.Marshal(errorHandler)
		if err != nil {
			return
		}

		writer.Write(buf)
	}
}

type ErrorWrapperHandler struct {
	Wrap http.Handler
}

func (h ErrorWrapperHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Print("preprocessing ServeHTTP")

	h.Wrap.ServeHTTP(writer, request)

	return
}
