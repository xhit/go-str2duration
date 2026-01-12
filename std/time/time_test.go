package time

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input string
		want  time.Duration
	}{
		{"1s", 1 * time.Second},
		{"500ms", 500 * time.Millisecond},
		{"1.5m", 90 * time.Second},
		{"1h30m", 90 * time.Minute},
		{"1d", 24 * time.Hour},
		{"1w", 7 * 24 * time.Hour},
		{"1s 500ms", 1500 * time.Millisecond}, // Space separator handling
		{"10us", 10 * time.Microsecond},
		{"10µs", 10 * time.Microsecond},
		{"10us45m2h15s", 10*time.Microsecond + 45*time.Minute + 2*time.Hour + 15*time.Second}, // Out-of-order time
	}

	for _, tt := range tests {
		got, err := ParseDuration(tt.input)
		if err != nil {
			t.Errorf("ParseDuration(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseDuration(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseDuration_Errors(t *testing.T) {
	invalidInputs := []string{
		"1kg",    // Wrong unit
		"hello",  // Garbage
		"",       // Empty
		"1.1.1s", // Bad number
	}

	for _, input := range invalidInputs {
		_, err := ParseDuration(input)
		if err == nil {
			t.Errorf("ParseDuration(%q) expected error, got nil", input)
		}
	}
}
