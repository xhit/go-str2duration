package parser

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/armourstill/str2quantity/unit"
)

// Number constrains the types that can be returned by Parse.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// safeSkipSeps skips allowed separators but preserves characters that start a valid number (digits, dot, signs).
func safeSkipSeps(s string, separators string) string {
	if separators == "" {
		// Default relaxed separators
		separators = " \t\n\r,;|/"
	}

	for len(s) > 0 {
		c := s[0]
		// Stop at number start (digits, dot, signs).
		if (c >= '0' && c <= '9') || c == '.' || c == '+' || c == '-' {
			return s
		}

		if strings.ContainsRune(separators, rune(c)) {
			s = s[1:]
			continue
		}

		// Unknown char found
		return s
	}
	return s
}

// Parse parses a string into a standardized numerical value and its dimension.
// It uses input unit.System for configuration.
//
// Constraints:
// 1. System base unit (Scale=1.0) must align with '1' of type N.
// 2. Fractional values in integer type N will return error.
func Parse[N Number](s string, sys *unit.System) (N, unit.Dimension, error) {
	// Epsilon handles floating point noise (e.g. for pico/nano prefixes).
	const epsilon = 1e-12

	var total N
	var detectedDim unit.Dimension
	isDimSet := false
	partsCount := 0

	orig := s

	// Initial skip
	s = safeSkipSeps(s, sys.Config.Separators)

	for s != "" {
		// Check multi-part restriction
		if partsCount > 0 && !sys.Config.AllowMultiPart {
			return 0, unit.Dimension{}, fmt.Errorf("multi-part format is not allowed for this unit system: %q", orig)
		}

		// 1. Parse number
		val, nextStr, err := parseNumber(s)
		if err != nil {
			return 0, unit.Dimension{}, err
		}
		s = nextStr

		// Skip separators between value and unit (e.g. "100 MB")
		s = safeSkipSeps(s, sys.Config.Separators)

		// 2. Parse unit string
		unitStr, nextStr := parseUnit(s, sys.Config.Separators)
		if unitStr == "" {
			return 0, unit.Dimension{}, fmt.Errorf("missing unit in %q", orig)
		}
		s = nextStr

		// 3. Resolve unit
		u, scaleRatio, found := sys.Resolve(unitStr)
		if !found {
			return 0, unit.Dimension{}, fmt.Errorf("unknown unit: %s", unitStr)
		}

		// 4. Dimension check
		if !isDimSet {
			detectedDim = u.Dimension
			isDimSet = true
		} else if !detectedDim.Equals(u.Dimension) {
			return 0, unit.Dimension{}, fmt.Errorf("mixed dimensions: %s and %s", detectedDim, u.Dimension)
		}

		// 5. Accumulate value (Value * PrefixScale * UnitScale)
		// Calculate the value in base units as float64 first.
		partVal := val * scaleRatio * u.Scale

		var partN N

		// Step A: Check if it's effectively an integer (handling float noise like 29.999995 -> 30).
		rounded := math.Round(partVal)
		if math.Abs(rounded-partVal) <= epsilon {
			// It is effectively an integer. Use the clean integer value to avoid truncating 29.999 to 29.
			partN = N(rounded)
		} else {
			// Step B: It is a "real" number with fractional part (e.g. 0.5 or 0.125).
			// Check if the target generic type N can represent it.
			castN := N(partVal)

			// If N is float64, castN should be equal to partVal (diff ~ 0).
			// If N is int64, castN will be truncated, so diff will be large.
			if math.Abs(float64(castN)-partVal) > epsilon {
				return 0, detectedDim, fmt.Errorf("precision loss: part value %g cannot be represented exactly in target type", partVal)
			}
			partN = castN
		}

		total += partN
		partsCount++

		// Loop end skip
		s = safeSkipSeps(s, sys.Config.Separators)
	}

	return total, detectedDim, nil
}

// parseNumber extracts a float number from the beginning of the string.
// Supports integers, floats, and scientific notation (e.g. 1.2, 1e5).
// TODO: Potentially return a flag indicating if the input was syntactically an integer (no dot, no negative exponent).
// This could guide stricter precision checks or optimizations downstream, distinguishing
// "1" (syntax integer) from "1.0" (syntax float) or "0.9999999999999999" (float noise).
func parseNumber(s string) (float64, string, error) {
	end := 0
	allowSign := true
	allowDot := true
	allowE := true

	for end < len(s) {
		c := s[end]
		if c >= '0' && c <= '9' {
			// digits are always ok
			allowSign = false
		} else if c == '.' && allowDot {
			allowDot = false
			allowSign = false
		} else if (c == 'e' || c == 'E') && allowE && end > 0 { // e must not be start
			allowE = false
			allowDot = false // no dots after e
			allowSign = true // sign allowed after e
		} else if (c == '+' || c == '-') && allowSign {
			allowSign = false
		} else {
			break
		}
		end++
	}

	if end == 0 {
		return 0, s, errors.New("invalid number")
	}

	val, err := strconv.ParseFloat(s[:end], 64)
	if err != nil {
		return 0, s, err
	}

	return val, s[end:], nil
}

// parseUnit extracts the unit string.
// It stops when it encounters a digit, various signs, or a configured separator.
func parseUnit(s string, separators string) (string, string) {
	if separators == "" {
		separators = " \t\n\r,;|/"
	}

	end := 0
	for end < len(s) {
		c := s[end]
		// Stop at digits, dot, plus, minus (start of next number)
		if unicode.IsDigit(rune(c)) || c == '.' || c == '+' || c == '-' {
			break
		}
		// Stop at separators
		if strings.ContainsRune(separators, rune(c)) {
			break
		}
		end++
	}
	return s[:end], s[end:]
}
