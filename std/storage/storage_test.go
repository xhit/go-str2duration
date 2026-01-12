package storage

import (
	"math"
	"testing"
)

func TestParseStorage(t *testing.T) {
	const k = 1024.0
	const m = k * 1024
	const g = m * 1024

	tests := []struct {
		input    string
		expected float64 // Expected Bytes
		hasError bool
	}{
		// Basic Bytes
		{"1B", 1, false},
		{"1Byte", 1, false},
		{"100 Bytes", 100, false},
		{"0B", 0, false},

		// Bits
		{"8b", 1, false},       // 8 bits = 1 Byte
		{"1bit", 0.125, false}, // 1 bit
		{"16bits", 2, false},   // 16 bits = 2 Bytes

		// SI Prefixes (Decimal 1000) - REPLACED by JEDEC Logic (Binary 1024)
		// Now lowercase 'k' is also 1024
		{"1kB", k, false},       // k = 1024
		{"1kbit", k / 8, false}, // 1024 bits = 128 Bytes
		{"1gB", g, false},       // Lowercase g = 1024^3, with B=Byte

		// JEDEC / Common Usage Override (Uppercase = 1024)
		{"1KB", k, false},     // K = 1024
		{"1MB", m, false},     // M = 1024^2
		{"1GB", g, false},     // G = 1024^3
		{"1Kb", k / 8, false}, // 1 Kbit = 1024 bits = 128 Bytes

		// IEC Prefixes (Binary 1024) - with case variants
		{"1KiB", k, false},
		{"1kiB", k, false}, // Lowercase ki
		{"1KIB", k, false}, // Uppercase KI
		{"1MiB", m, false},
		{"1miB", m, false}, // Lowercase mi
		{"1GiB", g, false},
		{"1Kib", k / 8, false}, // 1 Kibit = 1024 bits = 128 Bytes

		// Decimals
		{"1.5KB", 1.5 * k, false}, // 1.5 * 1024
		{"0.5B", 0.5, false},

		// Complex formatting
		{"  10 MB  ", 10 * m, false},

		// Errors
		{"10", 0, true},   // Missing unit
		{"10s", 0, true},  // Wrong dimension (time)
		{"10Kg", 0, true}, // Unknown unit
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		got, err := ParseBytes(tt.input)
		if tt.hasError {
			if err == nil {
				t.Errorf("ParseBytes(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseBytes(%q) unexpected error: %v", tt.input, err)
			}
			// Use a small epsilon for float comparison
			if math.Abs(got-tt.expected) > 0.0001 {
				t.Errorf("ParseBytes(%q) = %v, expected %v", tt.input, got, tt.expected)
			}
		}
	}
}

func TestCaseSensitivity(t *testing.T) {
	// 'b' is bit (0.125 Byte)
	val1, err := ParseBytes("1b")
	if err != nil {
		t.Fatal(err)
	}
	if val1 != 0.125 {
		t.Errorf("Expected 1b to be 0.125, got %v", val1)
	}
	// 'B' is Byte (1.0 Byte)
	val2, err := ParseBytes("1B")
	if err != nil {
		t.Fatal(err)
	}
	if val2 != 1.0 {
		t.Errorf("Expected 1B to be 1.0, got %v", val2)
	}
}

func TestParseBits(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		hasError bool
	}{
		// Basic
		{"1b", 1, false},
		{"1bit", 1, false},
		{"8b", 8, false},
		{"1B", 8, false}, // 1 Byte = 8 bits
		{"1KB", 1024 * 8, false},

		// Integer bit values from fractional bytes
		{"1.5B", 12, false}, // 1.5 * 8 = 12 bits

		// Fractional bits (should fail for int64)
		{"0.5b", 0, true},
		{"0.1B", 0, true}, // 0.8 bits

		// Large values (int64 limit checks)
		{"1 PiB", 1 << 53, false},
	}

	for _, tt := range tests {
		got, err := ParseBits(tt.input)
		if tt.hasError {
			if err == nil {
				t.Errorf("ParseBits(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseBits(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.expected {
				t.Errorf("ParseBits(%q) = %v, expected %v", tt.input, got, tt.expected)
			}
		}
	}
}
