package main

import (
	"fmt"
	"github.com/marcosalbano/platform-api/internal/handlers"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/hostname", handlers.HostnameHandler)
	http.HandleFunc("/memory", handlers.MemoryHandler)
	http.HandleFunc("/uptime", handlers.UptimeHandler)
	http.HandleFunc("/system", handlers.SystemHandler)
	http.HandleFunc("/network", handlers.NetworkHandler)
	http.HandleFunc("/filesystem", handlers.FilesystemHandler)
	http.Handle("/", handlers.DashboardHandler())
	fmt.Println("Platform API listening on :8081")

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}

}
