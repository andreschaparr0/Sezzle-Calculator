// Package handlers wires HTTP requests to the calculator engine. It owns
// request decoding, input validation, and mapping domain errors to HTTP
// status codes — the calculator package itself never touches HTTP.
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"

	"https://github.com/andreschaparr0/Sezzle-Calculator/backend/internal/calculator"
	"https://github.com/andreschaparr0/Sezzle-Calculator/backend/internal/httpx"
)

// CalculateRequest is the JSON body accepted by POST /api/calculate.
//
// B is a pointer so we can tell "omitted" apart from "explicitly 0", which
// matters for validation: every operation except Sqrt requires B to be
// present in the request.
type CalculateRequest struct {
	Operation string   `json:"operation"`
	A         *float64 `json:"a"`
	B         *float64 `json:"b,omitempty"`
}

// CalculateResponse is the JSON body returned on a successful calculation.
type CalculateResponse struct {
	Operation string   `json:"operation"`
	A         float64  `json:"a"`
	B         *float64 `json:"b,omitempty"`
	Result    float64  `json:"result"`
}

// Calculator exposes the HTTP handlers for the calculator API.
// It is a thin struct (currently stateless) so that it can grow to hold
// dependencies (e.g. a logger or metrics client) without changing the
// handler signatures used by main.go.
type Calculator struct{}

// NewCalculator constructs a Calculator handler set.
func NewCalculator() *Calculator {
	return &Calculator{}
}

// HandleCalculate handles POST /api/calculate.
//
// Request:  {"operation": "add", "a": 4, "b": 2}
// Response: {"operation": "add", "a": 4, "b": 2, "result": 6}
//
// For the unary "sqrt" operation, "b" may be omitted entirely:
// Request:  {"operation": "sqrt", "a": 16}
// Response: {"operation": "sqrt", "a": 16, "result": 4}
func (c *Calculator) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid JSON body: %v", err))
		return
	}

	op := calculator.Operation(req.Operation)

	if req.Operation == "" {
		httpx.WriteError(w, http.StatusBadRequest, "\"operation\" is required")
		return
	}
	if !op.IsValid() {
		httpx.WriteError(w, http.StatusBadRequest, fmt.Sprintf("unsupported operation %q; supported: %v", req.Operation, calculator.Operations()))
		return
	}
	if req.A == nil {
		httpx.WriteError(w, http.StatusBadRequest, "\"a\" is required")
		return
	}
	if op != calculator.Sqrt && req.B == nil {
		httpx.WriteError(w, http.StatusBadRequest, "\"b\" is required for operation "+req.Operation)
		return
	}
	// Defensive: standard JSON has no literal for NaN/Infinity, so the
	// decoder above already rejects such payloads before we get here.
	// These checks guard against any future decoder/library change and
	// keep the invariant ("a and b are always finite past this point")
	// explicit and enforced rather than assumed.
	if err := validateFinite("a", *req.A); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.B != nil {
		if err := validateFinite("b", *req.B); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	var b float64
	if req.B != nil {
		b = *req.B
	}

	result, err := calculator.Calculate(op, *req.A, b)
	if err != nil {
		status, message := mapCalculationError(err)
		httpx.WriteError(w, status, message)
		return
	}

	if err := validateFinite("result", result); err != nil {
		// The inputs were individually finite (validated above) but combined
		// into a non-representable result, e.g. very large exponents
		// overflowing to +Inf. This is a client input problem, not a server
		// bug, so it is reported as 422 rather than 500.
		httpx.WriteError(w, http.StatusUnprocessableEntity, "result is not a finite number: "+err.Error())
		return
	}

	httpx.WriteJSON(w, http.StatusOK, CalculateResponse{
		Operation: req.Operation,
		A:         *req.A,
		B:         req.B,
		Result:    result,
	})
}

// HandleOperations handles GET /api/operations, a small discovery endpoint
// listing every operation the service supports. Useful for the frontend to
// render available buttons without hardcoding the list, and for API
// consumers exploring the service.
func (c *Calculator) HandleOperations(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"operations": calculator.Operations(),
	})
}

// mapCalculationError maps a domain error from the calculator package to
// an HTTP status code and message. Domain errors (division by zero,
// negative sqrt) represent a semantically invalid — but well-formed —
// request, so they map to 422 Unprocessable Entity rather than 400.
func mapCalculationError(err error) (int, string) {
	switch {
	case errors.Is(err, calculator.ErrDivisionByZero):
		return http.StatusUnprocessableEntity, "division by zero is not allowed"
	case errors.Is(err, calculator.ErrNegativeSqrt):
		return http.StatusUnprocessableEntity, "cannot compute the square root of a negative number"
	case errors.Is(err, calculator.ErrUnsupportedOperation):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "internal error"
	}
}

// validateFinite rejects NaN and +/-Inf values, which JSON cannot represent
// and which would otherwise cause json.Marshal to fail silently downstream.
func validateFinite(field string, v float64) error {
	if math.IsNaN(v) {
		return fmt.Errorf("%q must be a finite number, got NaN", field)
	}
	if math.IsInf(v, 0) {
		return fmt.Errorf("%q must be a finite number, got Infinity", field)
	}
	return nil
}
