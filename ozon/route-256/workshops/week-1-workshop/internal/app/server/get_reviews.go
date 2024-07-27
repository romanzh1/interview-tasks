package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type GetReviewsReviewsResponse struct {
	Sku     int    `json:"sku"`
	Comment string `json:"comment"`
	UserID  int    `json:"user_id"`
}
type GetReviewsResponse struct {
	Reviews []GetReviewsReviewsResponse `json:"reviews"`
}

func (s *Server) GetReviews(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	idRaw := request.PathValue("id")
	sku, err := strconv.Atoi(idRaw)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("error"))
		return
	}

	reviews, err := s.reviewsService.GetReviewsBySku(request.Context(), sku)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("error"))
		return
	}

	response := GetReviewsResponse{}
	for _, r := range reviews {
		response.Reviews = append(response.Reviews, GetReviewsReviewsResponse{
			Sku:     r.Sku,
			Comment: r.Comment,
			UserID:  r.UserID,
		})
	}

	buf, err := json.Marshal(&response)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("error"))
		return
	}

	_, err = writer.Write(buf)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("error"))
		return
	}
}
