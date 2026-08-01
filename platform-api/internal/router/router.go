package router

import (
	"net/http"

	"github.com/marcosalbano/platform-api/internal/handlers"
	"github.com/marcosalbano/platform-api/internal/middleware"
)

// New registra todas las rutas de la aplicación y
// devuelve el router listo para utilizar.
func New(authMiddleware *middleware.AuthMiddleware) http.Handler {

	mux := http.NewServeMux()

	// -----------------------------
	// Public endpoints
	// -----------------------------

	mux.HandleFunc("/health", handlers.HealthHandler)

	// -----------------------------
	// Protected API endpoints
	// -----------------------------

	mux.Handle(
		"/hostname",
		authMiddleware.Handler(
			http.HandlerFunc(handlers.HostnameHandler),
		),
	)

	mux.Handle(
		"/memory",
		authMiddleware.Handler(
			http.HandlerFunc(handlers.MemoryHandler),
		),
	)

	mux.Handle(
		"/uptime",
		authMiddleware.Handler(
			http.HandlerFunc(handlers.UptimeHandler),
		),
	)

	mux.Handle(
		"/system",
		authMiddleware.Handler(
			http.HandlerFunc(handlers.SystemHandler),
		),
	)

	mux.Handle(
		"/network",
		authMiddleware.Handler(
			http.HandlerFunc(handlers.NetworkHandler),
		),
	)

	mux.Handle(
		"/filesystem",
		authMiddleware.Handler(
			http.HandlerFunc(handlers.FilesystemHandler),
		),
	)

	// -----------------------------
	// Dashboard
	// -----------------------------

	// El frontend es público.
	// Los datos que consume continúan protegidos por autenticación.
	mux.Handle("/", handlers.DashboardHandler())

	return mux
}
