package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.GET("/knockknock", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello bhai")
	})
	e.Logger.Print("Listening to port %s", port)
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
}
