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
		"/api/v1/hostname",
		authMiddleware.Handler(
			middleware.Authorize("system", "read")(
				http.HandlerFunc(handlers.HostnameHandler),
			),
		),
	)

	mux.Handle(
		"/api/v1/memory",
		authMiddleware.Handler(
			middleware.Authorize("system", "read")(
				http.HandlerFunc(handlers.MemoryHandler),
			),
		),
	)

	mux.Handle(
		"/api/v1/uptime",
		authMiddleware.Handler(
			middleware.Authorize("system", "read")(
				http.HandlerFunc(handlers.UptimeHandler),
			),
		),
	)

	mux.Handle(
		"/api/v1/system",
		authMiddleware.Handler(
			middleware.Authorize("system", "read")(
				http.HandlerFunc(handlers.SystemHandler),
			),
		),
	)

	mux.Handle(
		"/api/v1/network",
		authMiddleware.Handler(
			middleware.Authorize("network", "read")(
				http.HandlerFunc(handlers.NetworkHandler),
			),
		),
	)

	mux.Handle(
		"/api/v1/filesystem",
		authMiddleware.Handler(
			middleware.Authorize("filesystem", "read")(
				http.HandlerFunc(handlers.FilesystemHandler),
			),
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
