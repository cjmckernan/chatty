package socks

import (
	"github.com/gorilla/websocket"
	"log"
)

type Hub struct {
	// Registered clients
	clients map[*websocket.Conn]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *websocket.Conn

	// Unregister requests from clients
	unregister chan *websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		clients:    make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.clients[conn] = true
		case conn := <-h.unregister:
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				closeConn(conn)
			}
		case message := <-h.broadcast:
			for conn := range h.clients {
				// Ensure the connection is open and send the message
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error writing message to WebSocket client: %v", err)
					closeConn(conn)
					delete(h.clients, conn)
				}
			}
		}
	}
}

func closeConn(conn *websocket.Conn) {
	conn.Close()
}
