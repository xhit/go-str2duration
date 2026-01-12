package unit

import "fmt"

// Dimension represents the physical dimensions of a quantity.
// It uses the SI base quantities.
type Dimension struct {
	L int // Length (meter)
	M int // Mass (kilogram)
	T int // Time (second)
	I int // Electric current (ampere)
	K int // Thermodynamic temperature (kelvin)
	N int // Amount of substance (mole)
	J int // Luminous intensity (candela)

	// Extra allows for non-standard dimensions checking (e.g., "digital")
	Extra string
}

// Equals checks if two dimensions are identical.
func (d Dimension) Equals(other Dimension) bool {
	return d == other
}

// String returns a string representation of the dimension.
func (d Dimension) String() string {
	if d.Extra != "" {
		return fmt.Sprintf("Dim(%s)", d.Extra)
	}
	return fmt.Sprintf("L^%d M^%d T^%d I^%d K^%d N^%d J^%d", d.L, d.M, d.T, d.I, d.K, d.N, d.J)
}

// Common dimensions
var (
	DimDimensionless = Dimension{}
	DimTime          = Dimension{T: 1}
	DimLength        = Dimension{L: 1}
	DimMass          = Dimension{M: 1}
	DimTemp          = Dimension{K: 1}
	DimCurrent       = Dimension{I: 1}
	DimAmount        = Dimension{N: 1}
	DimLuminous      = Dimension{J: 1}
	DimStorage       = Dimension{Extra: "storage"}
)
