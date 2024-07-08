package loms

import (
	"route256/cart/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	host       string
	dialOption grpc.DialOption
	client     proto.LomsServiceClient
	conn       *grpc.ClientConn
}

func NewClient(host string) *Client {
	return &Client{
		host:       host,
		dialOption: grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

func (c *Client) Run() error {
	conn, err := grpc.NewClient(c.host, c.dialOption)
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
