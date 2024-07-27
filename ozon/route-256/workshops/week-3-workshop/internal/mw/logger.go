package mw

import (
	"context"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// для превращения protbuf структур в json используем google.golang.org/protobuf/encoding/protojson пакет а не encoding/json
	raw, _ := protojson.Marshal((req).(proto.Message))
	log.Printf("request: method: %v, req: %v\n", info.FullMethod, string(raw))

	if resp, err = handler(ctx, req); err != nil {
		log.Printf("response: method: %v, err: %v\n", info.FullMethod, err)
		return
	}

	rawResp, _ := protojson.Marshal((resp).(proto.Message))

	// Пример академический
	// Логировать все нельзя, в объекте могут быть чувствительные к безопасности данные (персональные)
	// Также сами объекты для логирования могут быть очень большими
	log.Printf("response: method: %v, resp: %v\n", info.FullMethod, string(rawResp))

	return
}

func WithHTTPLoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
