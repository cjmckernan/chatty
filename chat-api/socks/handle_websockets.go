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

	// Retrieve the session ID from the query parameters
	sessionID := c.QueryParam("sessionId")
	if sessionID == "" {
		log.Printf("Missing session ID in WebSocket request")
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing session ID")
	}

	// Get the username associated with the session ID
	username, err := message_store.GetUsernameBySessionID(sessionID)
	if err != nil {
		log.Printf("Invalid session ID: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid session ID")
	}

	hub.register <- ws
	hub.clients[ws] = username

	defer func() {
		hub.unregister <- ws
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			return err
		}

		log.Printf("Received message from %s: %s", username, string(msg))

		var message map[string]string
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("Error unmarshaling WebSocket message: %v", err)
			continue
		}

		// Validate the username in the message payload
		if message["username"] != username {
			log.Printf("Username mismatch: %s (message) != %s (session)", message["username"], username)
			
			ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "username mismatch"))
			hub.unregister <- ws
			closeConn(ws)
			break 
    }

		topic := message["topic"]
		text := message["text"]

		// Store the message in Redis
		if err := message_store.StoreMessage(topic, sessionID, username, text); err != nil {
			log.Printf("Error storing message: %v", err)
		}

		// Broadcast the received message to all clients
		hub.broadcast <- msg
	}

	return nil
}

