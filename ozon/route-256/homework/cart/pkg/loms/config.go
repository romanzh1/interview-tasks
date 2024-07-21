package loms

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"

	"route256/cart/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	host   string
	client proto.LomsServiceClient
	conn   *grpc.ClientConn
}

func NewClient(host string) *Client {
	return &Client{
		host: host,
	}
}

func (c *Client) Run() error {
	conn, err := grpc.NewClient(c.host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler(
			otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
			otelgrpc.WithPropagators(otel.GetTextMapPropagator()))))
	if err != nil {
		return err
	}

	c.conn = conn
	c.client = proto.NewLomsServiceClient(conn)

	return nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
