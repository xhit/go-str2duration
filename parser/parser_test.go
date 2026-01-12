package parser_test

import (
	"testing"

	"github.com/armourstill/str2quantity/parser"
	"github.com/armourstill/str2quantity/unit"
)

func createTestSystem() *unit.System {
	sys := unit.NewSystem(unit.SystemConfig{
		AllowMultiPart: true,
	})
	// Time units
	sys.Add("s", 1, unit.DimTime)
	sys.Add("m", 60, unit.DimTime)
	sys.Add("h", 3600, unit.DimTime)

	// Length units
	sys.Add("meter", 1, unit.DimLength)

	// Prefixes
	sys.AddPrefix("m", 0.001, "s", "meter") // milli
	return sys
}

func TestParse_Basic(t *testing.T) {
	sys := createTestSystem()

	tests := []struct {
		input   string
		wantVal float64
		wantErr bool
	}{
		{"1s", 1, false},
		{"1m", 60, false},
		{"1h30m", 5400, false}, // 3600 + 30*60
		{"1.5h", 5400, false},
		{"100ms", 0.1, false},    // 100 * 0.001
		{"", 0, false},           // Empty string -> 0
		{"1x", 0, true},          // Unknown unit
		{"1s 1meter", 0, true},   // Mixed dimension
		{"1m2ms", 60.002, false}, // Prefix with multi-part
		{"invalid", 0, true},     // No number
	}

	for _, tt := range tests {
		got, _, err := parser.Parse[float64](tt.input, sys)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if err == nil && got != tt.wantVal {
			t.Errorf("Parse(%q) = %g, want %g", tt.input, got, tt.wantVal)
		}
	}
}

func TestParse_Separators(t *testing.T) {
	// Configure system with custom separators
	sys := unit.NewSystem(unit.SystemConfig{
		AllowMultiPart: true,
		Separators:     ",|/ ",
	})
	sys.Add("d", 1, unit.DimLength) // dummy unit
	sys.Add("h", 1, unit.DimLength)

	tests := []struct {
		input   string
		wantVal float64
	}{
		{"1d,1h", 2},
		{"1d , 1h", 2},
		{"1d///1h", 2},
		{"1d|1h", 2},
		{"1d 1h", 2},
		// Tricky case: '3' is NOT a separator (implied by safeSkip logic),
		// but here we didn't add digit to separators in config, just testing normal separators.
		// If we had "3" in separators, "1h3d" should be parsed as 1h + 3d (3 is number), not 1h + d.
	}

	for _, tt := range tests {
		got, _, err := parser.Parse[float64](tt.input, sys)
		if err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
		}
		if got != tt.wantVal {
			t.Errorf("Parse(%q) = %g, want %g", tt.input, got, tt.wantVal)
		}
	}
}

func TestParse_TrickySeparators(t *testing.T) {
	// User set '3' as separator (bad practice, but testing robustness)
	sys := unit.NewSystem(unit.SystemConfig{
		AllowMultiPart: true,
		Separators:     "3 ",
	})
	sys.Add("h", 1, unit.DimTime)
	sys.Add("m", 1, unit.DimTime)

	// Input "1h3m"
	// 3 should be treated as number start for "3m", NOT skipped as separator.
	// So result should be 1 + 3 = 4.
	input := "1h3m"
	got, _, err := parser.Parse[float64](input, sys)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if got != 4 {
		t.Errorf("Parse(%q) with separator '3' = %g, want 4 (3 should be interpreted as number)", input, got)
	}
}

func TestParse_MultiPartRestriction(t *testing.T) {
	sys := unit.NewSystem(unit.SystemConfig{AllowMultiPart: false})
	sys.Add("B", 1, unit.DimStorage)

	// 1.5B is fine (single part, decimal)
	if _, _, err := parser.Parse[float64]("1.5B", sys); err != nil {
		t.Errorf("Single part failed: %v", err)
	}

	// 1B 2B should fail
	if _, _, err := parser.Parse[float64]("1B2B", sys); err == nil {
		t.Error("Multi part should fail but succeeded")
	}
}
