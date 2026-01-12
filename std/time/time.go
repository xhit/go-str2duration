package time

import (
	"errors"
	"time"

	"github.com/armourstill/str2quantity/parser"
	"github.com/armourstill/str2quantity/unit"
)

// System is the shared unit system for Time operations.
var System *unit.System

func init() {
	// Initialize system for Time strings (additive, case-sensitive).
	System = unit.NewSystem(unit.SystemConfig{
		AllowMultiPart:  true,
		CaseInsensitive: false, // Go duration strings are case sensitive (ms, not MS)
	})

	// Register Standard Units
	// Base: Nanosecond (ns) = 1.0 (aligns with time.Duration).

	// SI Time Units
	System.Add("ns", 1.0, unit.DimTime)
	System.Add("us", 1e3, unit.DimTime)
	System.Add("µs", 1e3, unit.DimTime) // Support micro symbol
	System.Add("ms", 1e6, unit.DimTime)
	System.Add("s", 1e9, unit.DimTime)

	// Common Time Units
	System.Add("m", 60*1e9, unit.DimTime)      // Minute
	System.Add("h", 3600*1e9, unit.DimTime)    // Hour
	System.Add("d", 24*3600*1e9, unit.DimTime) // Day
	System.Add("w", 604800*1e9, unit.DimTime)  // Week
}

// ParseDuration parses a duration string into time.Duration.
// Supports additive formats ("1h30m") and decimal values ("1.5h").
func ParseDuration(s string) (time.Duration, error) {
	val, dim, err := parser.Parse[time.Duration](s, System)
	if err != nil {
		return 0, err
	}

	// Validate Dimension
	if !dim.Equals(unit.DimTime) {
		return 0, errors.New("parsed quantity is not a time duration")
	}

	return val, nil
}
