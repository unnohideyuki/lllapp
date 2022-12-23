package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"lllapp-backend/router"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	g := e.Group("")
	router.Register(g)

	e.Logger.Fatal(e.Start(":8080"))
}
