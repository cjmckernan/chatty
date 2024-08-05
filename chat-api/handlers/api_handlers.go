package handlers

import (
	"chat-api/message_store"
	"chat-api/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type UserRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

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

func HandleMessages(c echo.Context) error {
	topic := c.Param("topic")

	messages, err := message_store.GetMessagesByTopic(topic)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to retrieve messages",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":  true,
		"messages": messages,
	})
}

func HelloWorldHandler(c echo.Context) error {

	response := map[string]string{
		"message": "Hello ",
	}
	return c.JSON(http.StatusOK, response)
}

func HandlerGetTopics(c echo.Context) error {
	topics, err := message_store.GetTopics()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to retrieve topics",
		})
	}

	// Return the topics as a JSON array
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"topics":  topics,
	})
}
