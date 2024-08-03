package handlers

import (
	"chat-api/message_store"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Create the user information
func HandleCreateUser(c echo.Context) error {
	user := new(message_store.User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request payload",
		})
	}

	message_store.CreateUser(*user)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"username":   user.Username,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// TODO Implement auth method

// TODO Update method to properly handle websocket sessions
// Handle the conection
func HandleWebsocketConn(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Error forming websocket connection")
		return err
	}
	defer ws.Close()
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message")
			return nil
		}
		log.Printf(string(msg))
		break
	}
	return nil
}
