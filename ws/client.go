package ws

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

// Client ssss
type Client struct {
	server *Server
	key    int
	conn   *websocket.Conn
}

// ListenRead shit
func (c *Client) ListenRead() {
	receivedtext := make([]byte, 100)
	for {
		n, err := c.conn.Read(receivedtext)
		if err == io.EOF {
			c.server.DeleteClient(c)
			return
		} else if err != nil {
			fmt.Printf("Received: %d bytes\n", n)
		}
	}
}

// Message shit
type Message struct {
	key  int
	body string
}
