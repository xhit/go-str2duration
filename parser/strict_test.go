package parser_test

import (
	"math"
	"testing"
	"time"

	"github.com/armourstill/str2quantity/parser"
	"github.com/armourstill/str2quantity/unit"
)

// createStrictIntSystem creates a system where base unit 'u' has Scale=1.
// Designed for testing integer parsing (N=int64).
func createStrictIntSystem() *unit.System {
	sys := unit.NewSystem(unit.SystemConfig{
		AllowMultiPart: true,
	})
	// Use arbitrary dimensions, e.g., Length for testing
	sys.Add("u", 1.0, unit.Dimension{L: 1})
	sys.Add("k", 1000.0, unit.Dimension{L: 1})
	return sys
}

func TestParse_Generics_Integers(t *testing.T) {
	sys := createStrictIntSystem()

	tests := []struct {
		name        string
		input       string
		wantVal     int64
		expectError bool
	}{
		// Valid cases
		{"Simple int", "10u", 10, false},
		{"Large int", "1000u", 1000, false},
		{"Decimal resulting in int", "1.5k", 1500, false}, // 1.5 * 1000 = 1500.0 -> 1500 (OK)
		{"MultiPart int", "1k 50u", 1050, false},

		// Precision Loss cases
		{"Decimal fraction of base unit", "0.5u", 0, true},       // 0.5 -> int(0) -> Loss
		{"Tiny decimal", "0.00001u", 0, true},                    // Loss
		{"Decimal resulting in non-int", "1.0005k", 0, true},     // 1000.5 -> Loss
		{"Complex loss", "1u 0.5u", 0, true},                     // Second part fails
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := parser.Parse[int64](tt.input, sys)
			if (err != nil) != tt.expectError {
				t.Errorf("Parse() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !tt.expectError && got != tt.wantVal {
				t.Errorf("Parse() = %v, want %v", got, tt.wantVal)
			}
		})
	}
}

func TestParse_Generics_Floats(t *testing.T) {
	sys := createStrictIntSystem()

	// Float parsing should NOT error on decimals
	tests := []struct {
		name    string
		input   string
		wantVal float64
	}{
		{"Fractional base", "0.5u", 0.5},
		{"Tiny fraction", "0.001u", 0.001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := parser.Parse[float64](tt.input, sys)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if math.Abs(got-tt.wantVal) > 1e-9 {
				t.Errorf("Parse() = %v, want %v", got, tt.wantVal)
			}
		})
	}
}

// Simulate actual Time Duration usage
func TestParse_TimeDuration_Strict(t *testing.T) {
	// Setup a Time system similar to the one in std/time
	sys := unit.NewSystem(unit.SystemConfig{AllowMultiPart: true})
	sys.Add("ns", 1.0, unit.DimTime)
	sys.Add("us", 1000.0, unit.DimTime)
	sys.Add("ms", 1000000.0, unit.DimTime)
	sys.Add("s", 1000000000.0, unit.DimTime)

	tests := []struct {
		name        string
		input       string
		wantVal     time.Duration
		expectError bool
	}{
		{"1ns", "1ns", 1, false},
		{"1us", "1us", 1000, false},
		{"0.5us", "0.5us", 500, false},          // 0.5 * 1000 = 500 (OK)
		{"0.5ns", "0.5ns", 0, true},             // Error: Cannot represent 0.5ns
		{"Truncated us", "1.000000001us", 0, true}, // 1000.000001 ns -> Error
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := parser.Parse[time.Duration](tt.input, sys)
			if (err != nil) != tt.expectError {
				t.Errorf("Parse() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !tt.expectError && got != tt.wantVal {
				t.Errorf("Parse() = %v, want %v", got, tt.wantVal)
			}
		})
	}
}
