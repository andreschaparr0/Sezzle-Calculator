// Command server runs the calculator microservice's HTTP API.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/andreschaparr0/Sezzle-Calculator/backend/internal/server"
)

func main() {
	port := getEnv("PORT", "8080")
	origins := parseOrigins(getEnv("ALLOWED_ORIGINS", "*"))

	router := server.NewRouter(server.Config{AllowedOrigins: origins})

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Run the server in a goroutine so the main goroutine can wait for a
	// shutdown signal
	go func() {
		log.Printf("calculator microservice listening on :%s (allowed origins: %v)", port, origins)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

// getEnv returns the value of the named environment variable, or fallback
// if it is unset or empty.
func getEnv(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}

// parseOrigins splits a comma-separated ALLOWED_ORIGINS value into a slice,
// trimming whitespace around each entry.
func parseOrigins(raw string) []string {
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			origins = append(origins, p)
		}
	}
	return origins
}
