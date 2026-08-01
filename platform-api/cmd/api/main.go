package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marcosalbano/platform-api/internal/auth"
	"github.com/marcosalbano/platform-api/internal/middleware"
	"github.com/marcosalbano/platform-api/internal/router"
)

func main() {

	// Authentication

	repository := auth.NewFileRepository("configs/users.json")

	authenticator := auth.NewBasicAuthenticator(repository)

	authMiddleware := middleware.NewAuthMiddleware(authenticator)

	// HTTP Router

	r := router.New(authMiddleware)

	fmt.Println("Platform API listening on :8081")

	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}
}
