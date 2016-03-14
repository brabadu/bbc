package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/brabadu/bbc/ws"

	"golang.org/x/net/websocket"
)

func startUDPServer(port int, server *ws.Server) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("UDP server listning on port " + strconv.Itoa(port))

	var buffer [1500]byte
	for {
		n, cliaddr, err := conn.ReadFromUDP(buffer[0:])
		if err != nil {
			panic(err)
		}
		msg := string(buffer[0:n])
		server.NewMessage(msg)
		fmt.Printf("Read '%s' from client %s\n", msg, cliaddr.String())
	}
}

func wsHandlerFactory(wsServer *ws.Server) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		wsServer.AddClient(ws)
	}
}

func main() {
	fmt.Println("hello bbc")
	rand.Seed(time.Now().UnixNano())
	serverPort := 1988

	wsServer := ws.NewServer()
	go startUDPServer(serverPort, wsServer)
	go wsServer.Listen()

	clientsPort := 8891
	http.HandleFunc("/debug", func(w http.ResponseWriter, req *http.Request) {
		s := websocket.Server{Handler: websocket.Handler(wsHandlerFactory(wsServer))}
		s.ServeHTTP(w, req)
	})
	fmt.Println("Websocket server listning on port " + strconv.Itoa(clientsPort))

	err := http.ListenAndServe(":"+strconv.Itoa(clientsPort), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
