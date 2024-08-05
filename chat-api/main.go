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

	e.GET("/", handlers.HelloWorldHandler, handlers.ValidateSessionID)

	e.GET("/ws", socks.HandleWebsocketConn, handlers.ValidateSessionID)

	e.POST("/auth", handlers.HandlAuth)
	e.POST("/user/create", handlers.HandleCreateUser)

	e.GET("/topics", handlers.HandlerGetTopics, handlers.ValidateSessionID)

	e.Logger.Fatal(e.Start(":9009"))
}
