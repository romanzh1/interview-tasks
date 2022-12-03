package client

const DefaultPort = 4242

// implement client design as in http package

// host string, ..... port int

// main
func main() {
	cli := client.New("127.0.0.1", client.DefaultPort)
}
