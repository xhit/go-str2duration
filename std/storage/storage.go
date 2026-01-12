package storage

import (
	"errors"

	"github.com/armourstill/str2quantity/parser"
	"github.com/armourstill/str2quantity/unit"
)

// System is the standard unit system for digital storage.
var System *unit.System

// bitsPerByte defines the conversion factor between Bits and Bytes.
const bitsPerByte = 8.0

func init() {
	// Initialize system: Check common usage (no multipart, correct case sensitivity).
	System = unit.NewSystem(unit.SystemConfig{
		AllowMultiPart:  false,
		CaseInsensitive: false,
	})

	// --- 1. Register Base Units ---
	// Bit (b) is base unit (Scale=1.0) for integer counting compatibility.

	// Bit (Base Unit)
	System.Add("b", 1.0, unit.DimStorage)
	System.Add("bit", 1.0, unit.DimStorage)
	System.Add("bits", 1.0, unit.DimStorage)

	// Byte (1 Byte = 8 bits)
	System.Add("B", bitsPerByte, unit.DimStorage)
	System.Add("Byte", bitsPerByte, unit.DimStorage)
	System.Add("Bytes", bitsPerByte, unit.DimStorage)
	// Duplicate unit removed for brevity

	targetUnits := []string{"B", "Byte", "Bytes", "b", "bit", "bits"}

	// --- 2. Register IEC Standard Prefixes (Binary 1024) ---
	// Explicitly register case variants so users can write "1kib" even in case-sensitive mode.
	const binaryKilo = 1024.0
	iecPrefixes := []struct {
		val  float64
		syms []string
	}{
		{float64(1 << 10), []string{"Ki", "ki", "KI"}}, // Ki = 2^10
		{float64(1 << 20), []string{"Mi", "mi", "MI"}}, // Mi = 2^20
		{float64(1 << 30), []string{"Gi", "gi", "GI"}}, // Gi = 2^30
		{float64(1 << 40), []string{"Ti", "ti", "TI"}}, // Ti = 2^40
		{float64(1 << 50), []string{"Pi", "pi", "PI"}}, // Pi = 2^50
		{float64(1 << 60), []string{"Ei", "ei", "EI"}}, // Ei = 2^60
	}
	for _, p := range iecPrefixes {
		for _, sym := range p.syms {
			System.AddPrefix(sym, p.val, targetUnits...)
		}
	}

	// --- 3. Register JEDEC/Binary Prefixes ---
	// Adopts JEDEC standard: K, M, G are 1024-based.
	// Maps both upper/lower case prefixes to binary scale for UX (overriding standard SI meaning of 'm', 'k').
	prefixes := []struct {
		sym string
		val float64
	}{
		// Kilo (2^10)
		{"k", float64(1 << 10)},
		{"K", float64(1 << 10)},
		// Mega (2^20)
		{"m", float64(1 << 20)},
		{"M", float64(1 << 20)},
		// Giga (2^30)
		{"g", float64(1 << 30)},
		{"G", float64(1 << 30)},
		// Tera (2^40)
		{"t", float64(1 << 40)},
		{"T", float64(1 << 40)},
		// Peta (2^50)
		{"p", float64(1 << 50)},
		{"P", float64(1 << 50)},
		// Exa (2^60)
		{"e", float64(1 << 60)},
		{"E", float64(1 << 60)},
	}
	for _, p := range prefixes {
		System.AddPrefix(p.sym, p.val, targetUnits...)
	}
}

// Bits parses a storage string and returns the exact quantity in bits.
// It uses int64 to enforce integer precision (rejecting fractional bits).
//
// LIMITATION: Since the underlying calculation is based on bits (int64),
// the maximum representable value is approx 1.15 Exabytes (2^63 bits).
// For larger values (e.g. Zettabytes), use ParseBytes which uses float64.
func ParseBits(s string) (int64, error) {
	valBits, dim, err := parser.Parse[int64](s, System)
	if err != nil {
		return 0, err
	}
	if !dim.Equals(unit.DimStorage) {
		return 0, errors.New("parsed quantity is not a storage unit")
	}
	return valBits, nil
}

// ParseBytes parses a storage string and returns the quantity in Bytes.
// It uses float64 internally to allow:
// 1. Handling values larger than 1 Exabyte (which exceeds int64 range when counted in bits).
// 2. Supporting fractional bytes (e.g. "4 bits" = 0.5 Bytes).
//
// Note: While this allows inputs that result in fractional bits (like "0.5 bit"),
// it prioritizes range and flexibility over strict physical bit validity.
func ParseBytes(s string) (float64, error) {
	// Parse as float64 bits first.
	valBits, dim, err := parser.Parse[float64](s, System)
	if err != nil {
		return 0, err
	}
	if !dim.Equals(unit.DimStorage) {
		return 0, errors.New("parsed quantity is not a storage unit")
	}
	// Convert bits to Bytes.
	return valBits / bitsPerByte, nil
}
