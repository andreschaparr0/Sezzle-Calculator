// Package server assembles the HTTP router for the calculator
// microservice: route registration and cross-cutting middleware (CORS,
// request logging). Kept separate from cmd/server so it can be exercised
// in tests via httptest without booting a real listener.
package server

import (
	"log"
	"net/http"
	"time"

	"github.com/andreschaparr0/Sezzle-Calculator/backend/internal/handlers"
)

// Config controls router behavior that varies between environments.
type Config struct {
	// AllowedOrigins lists the origins permitted by CORS for browser
	// requests (e.g. "http://localhost:5173" for a local Vite dev server).
	// Use []string{"*"} to allow any origin (fine for this take-home; a
	// production deployment would pin this down to known frontend hosts).
	AllowedOrigins []string
}

// NewRouter builds the complete http.Handler for the service: routes plus
// the CORS and logging middleware chain.
func NewRouter(cfg Config) http.Handler {
	mux := http.NewServeMux()

	calc := handlers.NewCalculator()

	mux.HandleFunc("GET /health", handlers.HandleHealth)
	mux.HandleFunc("POST /api/calculate", calc.HandleCalculate)
	mux.HandleFunc("GET /api/operations", calc.HandleOperations)

	var h http.Handler = mux
	h = withCORS(cfg.AllowedOrigins, h)
	h = withLogging(h)
	return h
}

// withLogging logs method, path, status code, and duration for every
// request — the minimum observability expected of a standalone service.
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}

// statusRecorder wraps http.ResponseWriter to capture the status code
// written by downstream handlers, purely for logging purposes.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// withCORS adds permissive-but-configurable CORS headers so a browser-based
// frontend served from a different origin (e.g. a local dev server, or a
// separately-deployed static site) can call this API directly.
func withCORS(allowedOrigins []string, next http.Handler) http.Handler {
	allowAll := len(allowedOrigins) == 0
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		if o == "*" {
			allowAll = true
		}
		allowed[o] = true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowAll {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" && allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
