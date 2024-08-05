package socks

import (
	"chat-api/message_store"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// WebSocket upgrader with default options
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Hub to manage all WebSocket connections
var hub = NewHub()

func init() {
	go hub.Run() // Start the Hub to manage WebSocket connections
}

func HandleWebsocketConn(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return err
	}
	defer ws.Close()

	hub.register <- ws

	defer func() {
		hub.unregister <- ws
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			return err
		}

		log.Printf("Received message: %s", string(msg))

		var message map[string]string
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("Error unmarshaling WebSocket message: %v", err)
			continue
		}

		topic := message["topic"]
		sessionID := message["sessionId"]
		username := message["username"]
		text := message["text"]

		// Store the message in Redis
		if err := message_store.StoreMessage(topic, sessionID, username, text); err != nil {
			log.Printf("Error storing message: %v", err)
		}

		// Broadcast the received message to all clients
		hub.broadcast <- msg
	}
}

