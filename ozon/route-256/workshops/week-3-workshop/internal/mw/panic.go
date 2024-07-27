package mw

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func Panic(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("panic: %v\n", e)
			err = status.Errorf(codes.Internal, "panic: %v", e)
		}
	}()
	resp, err = handler(ctx, req)
	return resp, err
}
