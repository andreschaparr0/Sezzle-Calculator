package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"https://github.com/andreschaparr0/Sezzle-Calculator/backend/internal/calculator"
)

// doCalculateRequest is a small test helper that POSTs the given raw JSON
// body to HandleCalculate and returns the recorded response.
func doCalculateRequest(t *testing.T, body string) *httptest.ResponseRecorder {
	t.Helper()
	c := NewCalculator()
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c.HandleCalculate(rec, req)
	return rec
}

func decodeCalculateResponse(t *testing.T, rec *httptest.ResponseRecorder) CalculateResponse {
	t.Helper()
	var resp CalculateResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body %q: %v", rec.Body.String(), err)
	}
	return resp
}

func decodeErrorResponse(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error body %q: %v", rec.Body.String(), err)
	}
	return resp.Error
}

func TestHandleCalculate_Success(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantResult float64
	}{
		{"add", `{"operation":"add","a":4,"b":2}`, 6},
		{"subtract", `{"operation":"subtract","a":10,"b":3}`, 7},
		{"multiply", `{"operation":"multiply","a":5,"b":6}`, 30},
		{"divide", `{"operation":"divide","a":10,"b":4}`, 2.5},
		{"power", `{"operation":"power","a":2,"b":8}`, 256},
		{"percentage", `{"operation":"percentage","a":10,"b":200}`, 20},
		{"sqrt without b", `{"operation":"sqrt","a":81}`, 9},
		{"sqrt with b ignored", `{"operation":"sqrt","a":81,"b":0}`, 9},
		{"negative operands", `{"operation":"add","a":-5,"b":-3}`, -8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doCalculateRequest(t, tt.body)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
			}
			resp := decodeCalculateResponse(t, rec)
			if resp.Result != tt.wantResult {
				t.Errorf("result = %v, want %v", resp.Result, tt.wantResult)
			}
			ct := rec.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}
		})
	}
}

func TestHandleCalculate_ValidationErrors(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"missing operation", `{"a":1,"b":2}`, http.StatusBadRequest},
		{"empty operation", `{"operation":"","a":1,"b":2}`, http.StatusBadRequest},
		{"unsupported operation", `{"operation":"modulo","a":1,"b":2}`, http.StatusBadRequest},
		{"missing a", `{"operation":"add","b":2}`, http.StatusBadRequest},
		{"missing b for binary op", `{"operation":"add","a":1}`, http.StatusBadRequest},
		{"malformed json", `{"operation":"add","a":1,`, http.StatusBadRequest},
		{"unknown field", `{"operation":"add","a":1,"b":2,"c":3}`, http.StatusBadRequest},
		{"non-numeric a", `{"operation":"add","a":"x","b":2}`, http.StatusBadRequest},
		{"infinity-producing input", `{"operation":"power","a":1e300,"b":1e300}`, http.StatusUnprocessableEntity},
		{"nan-producing input", `{"operation":"power","a":-8,"b":0.5}`, http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doCalculateRequest(t, tt.body)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if msg := decodeErrorResponse(t, rec); msg == "" {
				t.Error("expected non-empty error message")
			}
		})
	}
}

func TestHandleCalculate_DomainErrors(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"division by zero", `{"operation":"divide","a":10,"b":0}`},
		{"negative sqrt", `{"operation":"sqrt","a":-9}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doCalculateRequest(t, tt.body)
			if rec.Code != http.StatusUnprocessableEntity {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusUnprocessableEntity, rec.Body.String())
			}
			if msg := decodeErrorResponse(t, rec); msg == "" {
				t.Error("expected non-empty error message")
			}
		})
	}
}

// TestMapCalculationError exercises mapCalculationError directly, including
// the ErrUnsupportedOperation and unknown-error branches, which are not
// reachable through HandleCalculate because operation validity is already
// checked before calculator.Calculate is ever invoked.
func TestMapCalculationError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"division by zero", calculator.ErrDivisionByZero, http.StatusUnprocessableEntity},
		{"negative sqrt", calculator.ErrNegativeSqrt, http.StatusUnprocessableEntity},
		{"unsupported operation", calculator.ErrUnsupportedOperation, http.StatusBadRequest},
		{"unknown error", errors.New("boom"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, msg := mapCalculationError(tt.err)
			if status != tt.wantStatus {
				t.Errorf("status = %d, want %d", status, tt.wantStatus)
			}
			if msg == "" {
				t.Error("expected non-empty message")
			}
		})
	}
}

func TestHandleOperations(t *testing.T) {
	c := NewCalculator()
	req := httptest.NewRequest(http.MethodGet, "/api/operations", nil)
	rec := httptest.NewRecorder()
	c.HandleOperations(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp struct {
		Operations []string `json:"operations"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Operations) != 7 {
		t.Errorf("got %d operations, want 7: %v", len(resp.Operations), resp.Operations)
	}
}

func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	HandleHealth(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("status field = %q, want %q", resp.Status, "ok")
	}
}
