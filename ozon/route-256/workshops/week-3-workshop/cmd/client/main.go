package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/internal/client/notes"
	"gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/internal/model"

	desc "gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/pkg/api/notes/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := desc.NewNotesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	wrappedClient := notes.NewClient("user", client)

	id, err := wrappedClient.SaveNote(ctx, &model.Note{
		Title:   "hello",
		Content: "woooooooorld",
	})

	if err != nil {
		log.Fatalf("got error: %v", err)
	}
	fmt.Println(id)
}
