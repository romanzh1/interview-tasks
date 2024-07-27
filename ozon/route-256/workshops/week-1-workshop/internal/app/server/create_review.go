package server

import (
	"encoding/json"
	"errors"
	"gitlab.ozon.dev/week-1-workshop/internal/pkg/reviews/model"
	"io"
	"log"
	"net/http"
	"strconv"
)

type CreateReviewRequest struct {
	Sku     int    `json:"sku"`
	Comment string `json:"comment"`
	UserID  int    `json:"user_id"`
}

type CreateReviewResponse struct {
	Sku     int    `json:"sku"`
	Comment string `json:"comment"`
	UserID  int    `json:"user_id"`
}

func (s *Server) CreateReview(writer http.ResponseWriter, request *http.Request) error {
	writer.Header().Set("Content-Type", "application/json")

	idRaw := request.PathValue("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		return err
	}

	log.Printf("id=%d\n", id)

	buf, errRead := io.ReadAll(request.Body)
	defer request.Body.Close()

	if errRead != nil {
		return errRead
	}

	var r CreateReviewRequest
	err = json.Unmarshal(buf, &r)
	if err != nil {
		return err
	}

	if r.Sku < 1 {
		return errors.New("fail validation")
	}

	// validation
	review := model.Review{
		Sku:     r.Sku,
		Comment: r.Comment,
		UserID:  r.UserID,
	}

	createdReview, err := s.reviewsService.CreateReview(request.Context(), review)
	if err != nil {
		return err
	}

	response := CreateReviewResponse{
		Sku:     createdReview.Sku,
		Comment: createdReview.Comment,
		UserID:  createdReview.UserID,
	}

	responseBuf, err := json.Marshal(&response)
	if err != nil {
		return err
	}

	_, err = writer.Write(responseBuf)

	return nil
}
