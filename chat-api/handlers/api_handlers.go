package handlers

import (
	"chat-api/message_store"
	"chat-api/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type UserRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

// Create the user information
func HandleCreateUser(c echo.Context) error {
	userRequest := new(UserRequest)

	if err := c.Bind(userRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"success": "false",
			"error":   "invalid request payload",
		})
	}

	if userRequest.Password != userRequest.PasswordConfirm {
		return c.JSON(http.StatusConflict, map[string]string{
			"success": "false",
			"error":   "passwords do not match",
		})
	}

	userExists, err := message_store.UserExists(userRequest.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not check user existence",
		})
	}

	if userExists {
		return c.JSON(http.StatusConflict, map[string]string{
			"success": "false",
			"error":   "username already exists",
		})
	}

	hashedPass, err := utils.HashPassword(userRequest.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"success": "false",
			"error":   "error hashing password",
		})
	}

	user := message_store.User{
		Username:     userRequest.Username,
		PasswordHash: hashedPass,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

  sessionID, err := message_store.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "failed to create user",
			"success": "false",
		})
	}
  c.Response().Header().Set("X-Session-ID", sessionID) 
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success":    "true",
		"username":   user.Username,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// TODO Add check for existing sessionID
func HandlAuth(c echo.Context) error {
	userRequest := new(UserRequest)

	if err := c.Bind(userRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"success": "false",
			"error":   "invalid request payload",
		})
	}

	sessionID, err := message_store.Auth(userRequest.Username, userRequest.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"success": "false",
			"error":   err.Error(),
		})
	}
	fmt.Printf(sessionID)
	c.Response().Header().Set("X-Session-ID", sessionID)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": "true",
	})
}

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
