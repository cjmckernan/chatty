package handlers

import (
	"chat-api/message_store"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ValidateSessionID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if the session ID is passed as a query parameter (for WebSockets)
		sessionID := c.QueryParam("sessionId")

		// If it's not a WebSocket request, try to get it from the header
		if sessionID == "" {
			sessionID = c.Request().Header.Get("X-Session-ID")
		}

		if sessionID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"success": "false",
				"error":   "missing session ID",
			})
		}

		username, err := message_store.GetUsernameBySessionID(sessionID)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"success": "false",
				"error":   "invalid session ID",
			})
		}
		// Set the username in the context and in the response headers
		c.Set("username", username)
		c.Response().Header().Set("X-User-Name", username)
		return next(c)
	}
}
