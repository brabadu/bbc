package ws

import (
	"bytes"
	"log"
	"math/rand"
	"strings"

	"golang.org/x/net/websocket"
)

// Server the server struct
type Server struct {
	nextClientID uint64
	clients      map[string]*Client
	messages     chan Message
}

// NewServer shit
func NewServer() *Server {
	clients := make(map[string]*Client)
	messages := make(chan Message)
	return &Server{
		0,
		clients,
		messages,
	}
}

const keyArr = "01234567890abcdefghigklmnopkrstuvwxyzABCDEFGHIGKLMNOPKRSTUVWXYZ"
const keyArrLen = uint64(len(keyArr))

func genKey(source uint64) string {
	var buffer bytes.Buffer
	for source > 0 {
		buffer.WriteByte(keyArr[source%keyArrLen])
		source = source / keyArrLen
	}

	return buffer.String()
}

// AddClient shit
func (s *Server) AddClient(conn *websocket.Conn) {
	nextClientKeySource := s.nextClientID*1000 + uint64(rand.Intn(1000))
	nextClientKey := genKey(nextClientKeySource)

	s.nextClientID++
	client := NewClient(nextClientKey, conn, s)
	s.clients[nextClientKey] = client

	log.Printf("New client id: %s", nextClientKey)
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

	s.messages <- Message{
		splittedStr[0],
		splittedStr[1],
	}
}

// Listen shit
func (s *Server) Listen() {
	for {
		select {
		case m := <-s.messages:
			client, ok := s.clients[m.key]
			if !ok {
				log.Printf("Client [%s] dropped, message not delivered", m.key)
				continue
			}

			log.Printf("Sending to client [%s] message %s\n", m.key, m.body)
			_, err := client.conn.Write([]byte(m.body))
			if err != nil {
				log.Fatalf("Write error: %s", err)
			}
		}
	}
}
