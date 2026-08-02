package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/marcosalbano/platform-api/internal/auth"
	"github.com/marcosalbano/platform-api/internal/config"
	"github.com/marcosalbano/platform-api/internal/middleware"
	"github.com/marcosalbano/platform-api/internal/router"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	repository := auth.NewFileRepository(config.UsersFilePath())

	authenticator := auth.NewBasicAuthenticator(repository)

	authMiddleware := middleware.NewAuthMiddleware(authenticator)

	r := router.New(authMiddleware)

	fmt.Println("Platform API listening on :8081")

	server := &http.Server{
		Addr:              ":8081",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-stop:
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownContext); err != nil {
			return fmt.Errorf("gracefully shut down HTTP server: %w", err)
		}

		return nil
	}
}
