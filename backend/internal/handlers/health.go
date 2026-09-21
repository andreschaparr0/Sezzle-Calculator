package handlers

import (
	"net/http"

	"https://github.com/andreschaparr0/Sezzle-Calculator/backend/internal/httpx"
)

// HandleHealth handles GET /health. It's a standard liveness/readiness
// probe endpoint expected of any independently-deployable microservice
// (e.g. for container orchestrators or load balancers).
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
