package socks

import (
	"chat-api/message_store"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Hub struct {
	// Map WebSocket connections to their associated usernames
	clients map[*websocket.Conn]string

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
		clients:    make(map[*websocket.Conn]string), // Initialize map
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			// We don't set the username here because it will be added when the connection is registered
			h.clients[conn] = ""
		case conn := <-h.unregister:
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				closeConn(conn)
			}
		case message := <-h.broadcast:
			for conn, username := range h.clients {
				// Ensure the connection is open and send the message
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("Error writing message to WebSocket client (%s): %v", username, err)
					closeConn(conn)
					delete(h.clients, conn)
				}
			}
		}
	}
}

func (h *Hub) SubscribeToResidChannel(topic string) {
	pubsub := message_store.SubscribeToTopic(topic)
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var message message_store.Message
			if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
				log.Printf("Error processing message %w", err)
				continue
			}
			h.broadcast <- []byte(msg.Payload)
		}
	}()
}

func closeConn(conn *websocket.Conn) {
	conn.Close()
}
