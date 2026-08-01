package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marcosalbano/platform-api/internal/auth"
	"github.com/marcosalbano/platform-api/internal/handlers"
	"github.com/marcosalbano/platform-api/internal/middleware"
)

func main() {

	// --------------------------------------------------
	// Authentication subsystem
	// --------------------------------------------------

	repository := auth.NewFileRepository("configs/users.json")

	authenticator := auth.NewBasicAuthenticator(repository)

	authMiddleware := middleware.NewAuthMiddleware(authenticator)

	// --------------------------------------------------
	// Public endpoints
	// --------------------------------------------------

	http.HandleFunc("/health", handlers.HealthHandler)

	// --------------------------------------------------
	// Protected API endpoints
	// --------------------------------------------------

	http.Handle(
		"/hostname",
		authMiddleware.Handler(
			http.HandlerFunc(handlers.HostnameHandler),
		),
	)

	http.Handle(
		"/memory",
		authMiddleware.Handler(
			http.HandlerFunc(handlers.MemoryHandler),
		),
	)

	http.Handle(
		"/uptime",
		authMiddleware.Handler(
			http.HandlerFunc(handlers.UptimeHandler),
		),
	)

	http.Handle(
		"/system",
		authMiddleware.Handler(
			http.HandlerFunc(handlers.SystemHandler),
		),
	)

	http.Handle(
		"/network",
		authMiddleware.Handler(
			http.HandlerFunc(handlers.NetworkHandler),
		),
	)

	http.Handle(
		"/filesystem",
		authMiddleware.Handler(
			http.HandlerFunc(handlers.FilesystemHandler),
		),
	)

	// --------------------------------------------------
	// Dashboard
	// --------------------------------------------------

	http.Handle(
		"/",
		authMiddleware.Handler(
			handlers.DashboardHandler(),
		),
	)

	fmt.Println("Platform API listening on :8081")

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
