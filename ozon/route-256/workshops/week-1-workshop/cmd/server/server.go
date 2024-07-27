package main

import (
	"gitlab.ozon.dev/week-1-workshop/internal/app/server"
	"gitlab.ozon.dev/week-1-workshop/internal/pkg/middleware"
	"gitlab.ozon.dev/week-1-workshop/internal/pkg/reviews/repository"
	"gitlab.ozon.dev/week-1-workshop/internal/pkg/reviews/service"
	"log"
	"net"
	"net/http"
)

func main() {
	log.Print("starting app")

	conn, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()

	log.Print("starting listen")

	reviewsRepository := repository.NewRepository()
	reviewsService := service.NewService(reviewsRepository)

	reviewsServer := server.NewServer(reviewsService)

	log.Print("starting serve")

	mux := http.NewServeMux()
	mux.Handle("POST /products/{id}/reviews", middleware.ErrorWrapper(reviewsServer.CreateReview))
	mux.HandleFunc("GET /products/{id}/reviews", reviewsServer.GetReviews)

	h := middleware.ErrorWrapperHandler{
		Wrap: mux,
	}

	if err := http.Serve(conn, h); err != nil {
		log.Panic(err)
	}
}
