// starts up the server
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	bootstrap "github.com/erik-sostenes/auth-api/cmd/server/dependency"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	logger "github.com/labstack/gommon/log"
)

const defaultPort = "1324"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.SetFlags(log.Flags() | log.Lshortfile)

	e := echo.New()
	e.Logger.SetLevel(logger.INFO)

	e.Use(middleware.Recover(), middleware.Logger(), middleware.CORS())

	if err := bootstrap.Injector(e); err != nil {
		e.Logger.Fatal(err)
	}

	go func() {
		if err := e.Start(":" + defaultPort); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
