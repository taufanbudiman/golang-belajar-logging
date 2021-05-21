package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	middlewareLogging "parsing-http/middleware/library"
)

func middlewareOne(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("from middleware one")
		return next(c)
	}
}

func middlewareTwo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("from middleware two")
		return next(c)
	}
}

func middlewareSomething(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("from middleware something")
		next.ServeHTTP(w, r)
	})
}

func main() {
	e := echo.New()
	e.Use(middlewareOne)
	e.Use(middlewareTwo)
	e.Use(echo.WrapMiddleware(middlewareSomething))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middlewareLogging.MiddlewareLogging)
	e.HTTPErrorHandler = middlewareLogging.ErrorHandler
	// middleware here

	e.GET("/index", func(c echo.Context) (err error) {
		fmt.Println("threeeeee!")

		return c.JSON(http.StatusOK, true)
	})

	lock := make(chan error)
	go func(lock chan error) { lock <- e.Start(":8000") }(lock)

	time.Sleep(1 * time.Millisecond)
	middlewareLogging.MakeLogEntry(nil).Warning("application started without ssl/tls enabled")

	err := <-lock
	if err != nil {
		middlewareLogging.MakeLogEntry(nil).Panic("failed to start application")
	}
}
