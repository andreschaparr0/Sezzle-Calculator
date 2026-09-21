package calculator

import "errors"

// Sentinel errors returned by the calculator package.
// Handlers map these to appropriate HTTP status codes (see internal/handlers).
var (
	// ErrDivisionByZero is returned when a division or percentage operation
	// would divide by zero.
	ErrDivisionByZero = errors.New("division by zero")

	// ErrNegativeSqrt is returned when attempting to take the square root
	// of a negative number (unsupported: no complex number support).
	ErrNegativeSqrt = errors.New("cannot compute square root of a negative number")

	// ErrUnsupportedOperation is returned when the requested operation is
	// not recognized by the calculator engine.
	ErrUnsupportedOperation = errors.New("unsupported operation")
)
