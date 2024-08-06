package main

import (
	"chat-api/handlers"
	"chat-api/socks"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "X-Session-ID"},
		ExposeHeaders: []string{"X-Session-ID", "X-User-Name"},
	}))

	e.GET("/ping", handlers.Ping, handlers.ValidateSessionID)

	//Websockets endpoints
	e.GET("/ws", socks.HandleWebsocketConn, handlers.ValidateSessionID)

	// User Endpoints
	e.POST("/auth", handlers.HandlAuth)
	e.POST("/user/create", handlers.HandleCreateUser)

	// Topics
	e.GET("/topics", handlers.HandlerGetTopics, handlers.ValidateSessionID)

	// Messages endpoints
	e.GET("/messages/:topic", handlers.HandleMessages, handlers.ValidateSessionID)

	e.Logger.Fatal(e.Start(":9009"))
}
