package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_HealthEndpoint(t *testing.T) {
	router := NewRouter(Config{AllowedOrigins: []string{"*"}})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestRouter_CalculateEndpoint(t *testing.T) {
	router := NewRouter(Config{AllowedOrigins: []string{"*"}})
	body := `{"operation":"add","a":2,"b":3}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp struct {
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Result != 5 {
		t.Errorf("result = %v, want 5", resp.Result)
	}
}

func TestRouter_OperationsEndpoint(t *testing.T) {
	router := NewRouter(Config{AllowedOrigins: []string{"*"}})
	req := httptest.NewRequest(http.MethodGet, "/api/operations", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestRouter_UnknownRoute(t *testing.T) {
	router := NewRouter(Config{AllowedOrigins: []string{"*"}})
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestRouter_MethodNotAllowed(t *testing.T) {
	router := NewRouter(Config{AllowedOrigins: []string{"*"}})
	// /api/calculate only accepts POST.
	req := httptest.NewRequest(http.MethodGet, "/api/calculate", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestWithCORS_AllowAllOrigins(t *testing.T) {
	router := NewRouter(Config{AllowedOrigins: []string{"*"}})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "http://example.com")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "*")
	}
}

func TestWithCORS_AllowSpecificOrigin(t *testing.T) {
	router := NewRouter(Config{AllowedOrigins: []string{"http://localhost:5173"}})

	t.Run("allowed origin echoed back", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
			t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:5173")
		}
	})

	t.Run("disallowed origin not echoed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		req.Header.Set("Origin", "http://evil.example.com")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Errorf("Access-Control-Allow-Origin = %q, want empty", got)
		}
	})
}

func TestWithCORS_PreflightRequest(t *testing.T) {
	router := NewRouter(Config{AllowedOrigins: []string{"*"}})
	req := httptest.NewRequest(http.MethodOptions, "/api/calculate", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestLoggingMiddleware_CapturesStatus(t *testing.T) {
	router := NewRouter(Config{})
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	rec := httptest.NewRecorder()

	// This mainly exercises withLogging + statusRecorder for coverage;
	// correctness of the log line itself isn't asserted since it writes
	// to the standard logger rather than being returned.
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
