package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	notes_usecase "gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/internal/service/notes"

	notes_repository "gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/internal/repository/notes"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc/credentials/insecure"

	"gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/internal/app/notes"
	"gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/internal/mw"
	desc "gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/pkg/api/notes/v1"
	"google.golang.org/grpc/reflection"

	"log"
	"net"

	"google.golang.org/grpc"
)

/*
*
Порты, которые слушает сервер
*/
const (
	grpcPort = 50051
	httpPort = 8081
)

func headerMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "x-auth":
		return key, true
	default:
		return key, false
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mw.Panic,
			mw.Logger,
			mw.Auth,
			mw.Validate,
		),
	)

	// Рефлексия - это возможность диктовать клиентам свой контракт
	reflection.Register(grpcServer)

	repository := notes_repository.NewRepository()
	useCase := notes_usecase.NewService(repository)
	controller := notes.NewService(useCase)

	// Сгенерированный метод из прото
	desc.RegisterNotesServer(grpcServer, controller)

	log.Printf("server listening at %v", lis.Addr())

	// grpcServer.Serve блокирующий, поэтому создаем в отдельной рутине
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// создаем клиента на наш grpc сервер
	conn, err := grpc.Dial(":50051", grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to deal:", err)
	}

	// создаем http прокси над нашим grpc
	// gwmux := runtime.NewServeMux()

	//gwmux.Handle("/swaggerui", )

	gwmux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(headerMatcher))

	if err = desc.RegisterNotesHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: mw.WithHTTPLoggingMiddleware(gwmux),
	}

	log.Printf("Serving gRPC-Gateway on %s\n", gwServer.Addr)

	log.Fatalln(gwServer.ListenAndServe())
}
