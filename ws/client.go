package ws

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

// Client ssss
type Client struct {
	key    string
	conn   *websocket.Conn
	server *Server
}

// NewClient shit
func NewClient(key string, connection *websocket.Conn, server *Server) *Client {
	client := &Client{
		key,
		connection,
		server,
	}
	client.conn.Write([]byte(key))
	return client
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
	key  string
	body string
}
