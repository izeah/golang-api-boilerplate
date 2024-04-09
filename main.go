package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"boilerplate/internal/config"
	"boilerplate/internal/delivery"
	"boilerplate/internal/factory"
	"boilerplate/internal/middleware"
	"boilerplate/pkg/database"
	"boilerplate/pkg/redis"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func init() {
	_ = godotenv.Load("./deploy/dev.env")
	logrus.Info("Choosen Environment : ", os.Getenv("ENV"))
}

// @title BAF QRCode API CMS
// @version 0.1.0
// @description This is a doc for baf-qrcode-api-cms.

// host localhost:4040
// BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	PORT := fmt.Sprintf("%d", config.App().Port)

	database.Init()
	defer database.Close()

	redis.Init()
	defer redis.Close()

	e := echo.New()
	f := factory.NewFactory()

	middleware.Init(e)
	delivery.HTTP(e, f)

	// Start server
	go func() {
		if err := e.Start(":" + PORT); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	// Recieve shutdown signals.
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
