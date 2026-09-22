// Package calculator contains the pure arithmetic engine used by the
// calculator microservice. It has zero knowledge of HTTP, JSON, or any
// transport concern — that separation keeps the business logic trivial
// to unit test and reusable from any future transport (REST, gRPC, CLI...).
package calculator

import "math"

// Operation identifies an arithmetic operation supported by the engine.
type Operation string

// Supported operations. Add/Subtract/Multiply/Divide are the required
// "basic" operations; Power/Sqrt/Percentage are the "advanced" ones called
// out as optional in the  assignment
const (
	Add        Operation = "add"
	Subtract   Operation = "subtract"
	Multiply   Operation = "multiply"
	Divide     Operation = "divide"
	Power      Operation = "power"
	Sqrt       Operation = "sqrt"
	Percentage Operation = "percentage"
)

// binaryOps holds every operation that takes two operands (a, b).
// Sqrt is unary and is handled separately in Calculate.
var binaryOps = map[Operation]func(a, b float64) (float64, error){
	Add:        func(a, b float64) (float64, error) { return a + b, nil },
	Subtract:   func(a, b float64) (float64, error) { return a - b, nil },
	Multiply:   func(a, b float64) (float64, error) { return a * b, nil },
	Divide:     divide,
	Power:      power,
	Percentage: percentage,
}

// Calculate dispatches to the requested Operation and returns its result.
//
// For the unary Sqrt operation, only a is used; b is ignored.
// Returns ErrUnsupportedOperation if op is not recognized.
func Calculate(op Operation, a, b float64) (float64, error) {
	if op == Sqrt {
		return sqrt(a)
	}

	fn, ok := binaryOps[op]
	if !ok {
		return 0, ErrUnsupportedOperation
	}
	return fn(a, b)
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

func power(base, exponent float64) (float64, error) {
	return math.Pow(base, exponent), nil
}

func sqrt(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSqrt
	}
	return math.Sqrt(a), nil
}

// percentage computes "a percent of b" as (a * b) / 100, e.g.
// percentage(10, 200) == 20 (10% of 200). This mirrors the most common
// calculator UX for the "%" operation. b is not a divisor here, so no
// division-by-zero case exists (b == 0 simply yields 0).
func percentage(a, b float64) (float64, error) {
	return (a * b) / 100, nil
}

// Operations returns the list of all supported operations, useful for
// input validation and for advertising capabilities (e.g. a discovery
// endpoint) without hardcoding the list in multiple places.
func Operations() []Operation {
	return []Operation{Add, Subtract, Multiply, Divide, Power, Sqrt, Percentage}
}

// IsValid reports whether op is a recognized Operation.
func (op Operation) IsValid() bool {
	if op == Sqrt {
		return true
	}
	_, ok := binaryOps[op]
	return ok
}
