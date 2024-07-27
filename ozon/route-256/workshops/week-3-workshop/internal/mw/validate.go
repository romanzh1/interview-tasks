package mw

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
Validate

	Любой интерцептор принимает:
	ctx - уже знакомый нам контекст
	req - может быть абсолютно любым, так как разные хендлеры принимают разный запрос
	info - информация о методе и сервере, который мы вызываем
	handler - то, что будет вызвано после работы интерцептора
*/
func Validate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if v, ok := req.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return handler(ctx, req)
}
