package unit_test

import (
	"testing"

	"github.com/armourstill/str2quantity/unit"
)

func TestSystem_Resolve(t *testing.T) {
	sys := unit.NewSystem(unit.SystemConfig{CaseInsensitive: false})

	// Setup units
	sys.Add("m", 1.0, unit.DimLength)
	sys.Add("g", 1.0, unit.DimMass)

	// Setup prefixes
	// k binds to m and g
	if err := sys.AddPrefix("k", 1000, "m", "g"); err != nil {
		t.Fatalf("failed to add prefix: %v", err)
	}
	// m (milli) only binds to m (meter)
	if err := sys.AddPrefix("m", 0.001, "m"); err != nil {
		t.Fatalf("failed to add prefix: %v", err)
	}

	tests := []struct {
		input          string
		wantScale      float64 // Total scale = prefixScale * unitScale
		wantUnitSymbol string
		found          bool
	}{
		{"m", 1.0, "m", true},           // Exact match
		{"km", 1000.0, "m", true},       // Prefix k + m
		{"mm", 0.001, "m", true},        // Prefix m + m
		{"kg", 1000.0, "g", true},       // Prefix k + g
		{"mg", 0, "", false},            // Prefix m + g (NOT ALLOWED by binding)
		{"x", 0, "", false},             // Unknown
		{"kx", 0, "", false},            // Prefix k + Unknown
	}

	for _, tt := range tests {
		u, prefixScale, found := sys.Resolve(tt.input)
		if found != tt.found {
			t.Errorf("Resolve(%q) found = %v, same as want? %v", tt.input, found, tt.found)
			continue
		}
		if !found {
			continue
		}
		if u.Symbol != tt.wantUnitSymbol {
			t.Errorf("Resolve(%q) unit = %s, want %s", tt.input, u.Symbol, tt.wantUnitSymbol)
		}
		totalScale := prefixScale * u.Scale
		if totalScale != tt.wantScale {
			t.Errorf("Resolve(%q) scale = %g, want %g", tt.input, totalScale, tt.wantScale)
		}
	}
}

func TestSystem_CloneAndOverwrite(t *testing.T) {
	sys := unit.NewSystem(unit.SystemConfig{})
	sys.Add("B", 1.0, unit.DimStorage)
	sys.AddPrefix("k", 1000, "B")

	// Create clone
	binarySys := sys.Clone()
	if err := binarySys.OverwritePrefix("k", 1024); err != nil {
		t.Fatalf("failed to overwrite prefix: %v", err)
	}

	// Check original intact
	_, scale1, _ := sys.Resolve("kB")
	if scale1 != 1000 {
		t.Errorf("Original system modified! k=%g, want 1000", scale1)
	}

	// Check clone modified
	_, scale2, _ := binarySys.Resolve("kB")
	if scale2 != 1024 {
		t.Errorf("Clone system not modified! k=%g, want 1024", scale2)
	}
}

func TestSystem_CaseInsensitive(t *testing.T) {
	sys := unit.NewSystem(unit.SystemConfig{CaseInsensitive: true})
	sys.Add("m", 1.0, unit.DimLength)
	sys.AddPrefix("k", 1000, "m")

	tests := []string{"m", "M", "km", "KM", "Km", "kM"}
	for _, input := range tests {
		_, _, found := sys.Resolve(input)
		if !found {
			t.Errorf("Resolve(%q) failed in case-insensitive mode", input)
		}
	}
}
