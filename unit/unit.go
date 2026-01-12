package unit

// Unit represents a measurement unit.
type Unit struct {
	Symbol    string
	Dimension Dimension
	Scale     float64 // Scale relative to the base unit of the dimension (e.g. 1000 for km if base is m)
}

// Prefix represents a unit prefix (e.g., "k" for kilo, "m" for milli).
type Prefix struct {
	Symbol string
	Scale  float64
}
