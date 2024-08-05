package socks

import (
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

// HandleWebsocketConn handles WebSocket connections
func HandleWebsocketConn(c echo.Context) error {
	// Upgrade HTTP to WebSocket protocol
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return err
	}
	defer ws.Close()

	hub.register <- ws // Register the new client

	defer func() {
		hub.unregister <- ws // Unregister when connection is closed
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			return err
		}

		log.Printf("Received message: %s", string(msg))

		// Broadcast the received message to all clients
		hub.broadcast <- msg
	}
}

