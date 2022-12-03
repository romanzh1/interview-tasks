package client

import "time"

const DefaultPort = 4242

// implement client design as in http package

// host string, ..... port int

// main
func main() {
	cli := client.New("127.0.0.1", client.DefaultPort)
}

// TODO resolve, add default client

type Client struct {
	host    string
	port    int
	timeout time.Time
}

var DefalutTimeout = time.Minutes * 10

func New(host string, port int) *Client {

	return &Client{
		host:    host,
		port:    port,
		timeout: timeout,
	}
}

var client = Client{}

// functional options

func New2(c Client) {

}
