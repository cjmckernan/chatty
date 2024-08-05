package main

import (
	"chat-api/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

  e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "X-Session-ID"},
    ExposeHeaders: []string{"X-Session-ID"}, 
  }))

	e.GET("/", func(c echo.Context) error {
		response := map[string]string{
			"message": "Hello World!",
		}
		return c.JSON(http.StatusOK, response)
	})

	e.GET("/ws", handlers.HandleWebsocketConn)

	e.POST("/auth", handlers.HandlAuth)
	e.POST("/user/create", handlers.HandleCreateUser)

	e.Logger.Fatal(e.Start(":9009"))
}
