package calculator

import (
	"errors"
	"math"
	"testing"
)

const epsilon = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestCalculate_BasicOperations(t *testing.T) {
	tests := []struct {
		name    string
		op      Operation
		a, b    float64
		want    float64
		wantErr error
	}{
		{"add positive", Add, 2, 3, 5, nil},
		{"add negative", Add, -2, -3, -5, nil},
		{"add floats", Add, 1.5, 2.25, 3.75, nil},
		{"subtract", Subtract, 10, 4, 6, nil},
		{"subtract negative result", Subtract, 4, 10, -6, nil},
		{"multiply", Multiply, 6, 7, 42, nil},
		{"multiply by zero", Multiply, 6, 0, 0, nil},
		{"multiply negative", Multiply, -3, 4, -12, nil},
		{"divide", Divide, 10, 2, 5, nil},
		{"divide with remainder", Divide, 7, 2, 3.5, nil},
		{"divide by zero", Divide, 5, 0, 0, ErrDivisionByZero},
		{"divide zero numerator", Divide, 0, 5, 0, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.op, tt.a, tt.b)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Calculate(%v, %v, %v) error = %v, want %v", tt.op, tt.a, tt.b, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Calculate(%v, %v, %v) unexpected error: %v", tt.op, tt.a, tt.b, err)
			}
			if !almostEqual(got, tt.want) {
				t.Errorf("Calculate(%v, %v, %v) = %v, want %v", tt.op, tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestCalculate_AdvancedOperations(t *testing.T) {
	tests := []struct {
		name    string
		op      Operation
		a, b    float64
		want    float64
		wantErr error
	}{
		{"power positive exponent", Power, 2, 10, 1024, nil},
		{"power zero exponent", Power, 5, 0, 1, nil},
		{"power negative exponent", Power, 2, -1, 0.5, nil},
		{"power fractional base", Power, 2.5, 2, 6.25, nil},
		{"sqrt perfect square", Sqrt, 16, 0, 4, nil},
		{"sqrt zero", Sqrt, 0, 0, 0, nil},
		{"sqrt non-perfect square", Sqrt, 2, 0, math.Sqrt2, nil},
		{"sqrt negative", Sqrt, -4, 0, 0, ErrNegativeSqrt},
		{"percentage of value", Percentage, 10, 200, 20, nil},
		{"percentage 100", Percentage, 100, 50, 50, nil},
		{"percentage of zero base", Percentage, 10, 0, 0, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.op, tt.a, tt.b)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Calculate(%v, %v, %v) error = %v, want %v", tt.op, tt.a, tt.b, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Calculate(%v, %v, %v) unexpected error: %v", tt.op, tt.a, tt.b, err)
			}
			if !almostEqual(got, tt.want) {
				t.Errorf("Calculate(%v, %v, %v) = %v, want %v", tt.op, tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestCalculate_UnsupportedOperation(t *testing.T) {
	_, err := Calculate(Operation("modulo"), 1, 2)
	if !errors.Is(err, ErrUnsupportedOperation) {
		t.Fatalf("Calculate with unknown op error = %v, want %v", err, ErrUnsupportedOperation)
	}
}

func TestOperation_IsValid(t *testing.T) {
	for _, op := range Operations() {
		if !op.IsValid() {
			t.Errorf("Operations() returned %v but IsValid() is false", op)
		}
	}

	if Operation("bogus").IsValid() {
		t.Error("expected bogus operation to be invalid")
	}
}

func TestOperations_ContainsAllExpected(t *testing.T) {
	want := map[Operation]bool{
		Add: true, Subtract: true, Multiply: true, Divide: true,
		Power: true, Sqrt: true, Percentage: true,
	}
	got := Operations()
	if len(got) != len(want) {
		t.Fatalf("Operations() returned %d ops, want %d", len(got), len(want))
	}
	for _, op := range got {
		if !want[op] {
			t.Errorf("Operations() returned unexpected op %v", op)
		}
	}
}
