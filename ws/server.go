package ws

import (
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
)

// Server the server struct
type Server struct {
	nextClientKey int
	clients       map[int]*Client
	messages      chan Message
}

// NewServer shit
func NewServer() *Server {
	clients := make(map[int]*Client)
	messages := make(chan Message)
	return &Server{
		0,
		clients,
		messages,
	}
}

// AddClient shit
func (s *Server) AddClient(conn *websocket.Conn) {
	nextClientKey := s.nextClientKey
	s.nextClientKey++
	client := &Client{
		s,
		nextClientKey,
		conn,
	}
	client.conn.Write([]byte(strconv.Itoa(nextClientKey)))
	s.clients[nextClientKey] = client

	log.Printf("New client id: %d", nextClientKey)
	client.ListenRead()
}

// DeleteClient shit
func (s *Server) DeleteClient(c *Client) {
	delete(s.clients, c.key)
	log.Printf("Client connection dropped. %d clients left", len(s.clients))
}

// NewMessage shit
func (s *Server) NewMessage(str string) {
	splittedStr := strings.SplitN(str, "|", 2)

	id, err := strconv.Atoi(splittedStr[0])
	if err != nil {
		return
	}

	s.messages <- Message{
		id,
		splittedStr[1],
	}
}

// Listen shit
func (s *Server) Listen() {
	for {
		select {
		case c := <-s.messages:
			log.Printf("Sending to client [%d] message %s", c.key, c.body)
			_, err := s.clients[c.key].conn.Write([]byte(c.body))
			if err != nil {
				log.Fatalf("Write error: %s", err)
			}
		}
	}
}
