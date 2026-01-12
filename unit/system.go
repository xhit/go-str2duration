package unit

import (
	"fmt"
	"sort"
	"strings"
)

// SystemConfig configures the behavior of the unit system.
type SystemConfig struct {
	// AllowMultiPart enables summing multiple parts (e.g. "1h30m").
	// If false, only single part allowed (e.g. "1.5MB").
	AllowMultiPart bool

	// CaseInsensitive normalizes input to lowercase.
	CaseInsensitive bool

	// Separators allowed between parts (ignored during parsing).
	// Defaults to " \t\n\r,;|/" if empty.
	Separators string
}

// System is a registry for units and prefixes.
type System struct {
	units        map[string]Unit
	prefixes     []Prefix
	Config       SystemConfig

	// unitPrefixes maps unit symbol -> allowed prefix symbols.
	unitPrefixes map[string]map[string]bool
}

// NewSystem creates a new unit system with the given configuration.
func NewSystem(config SystemConfig) *System {
	return &System{
		units:        make(map[string]Unit),
		prefixes:     make([]Prefix, 0),
		unitPrefixes: make(map[string]map[string]bool),
		Config:       config,
	}
}

// normalizeKey adjusts the key based on case sensitivity settings.
func (s *System) normalizeKey(k string) string {
	if s.Config.CaseInsensitive {
		return strings.ToLower(k)
	}
	return k
}

// Add registers a new unit.
func (s *System) Add(symbol string, scale float64, dim Dimension) {
	key := s.normalizeKey(symbol)
	s.units[key] = Unit{Symbol: symbol, Scale: scale, Dimension: dim}
}

// AddPrefix registers a new prefix and binds it to specific units.
func (s *System) AddPrefix(prefixSymbol string, scale float64, targetUnits ...string) error {
	pKey := s.normalizeKey(prefixSymbol)

	// 1. Register or update prefix definition
	exists := false
	for _, p := range s.prefixes {
		if p.Symbol == pKey {
			if p.Scale != scale {
				return fmt.Errorf("prefix %s already defined with different scale", prefixSymbol)
			}
			exists = true
			break
		}
	}
	if !exists {
		s.prefixes = append(s.prefixes, Prefix{Symbol: pKey, Scale: scale})
		// Sort prefixes by length (longest first)
		sort.Slice(s.prefixes, func(i, j int) bool {
			return len(s.prefixes[i].Symbol) > len(s.prefixes[j].Symbol)
		})
	}

	// 2. Bind to target units
	for _, uSymbol := range targetUnits {
		uKey := s.normalizeKey(uSymbol)

		if _, ok := s.units[uKey]; !ok {
			return fmt.Errorf("cannot bind prefix to unknown unit: %s", uSymbol)
		}

		if s.unitPrefixes[uKey] == nil {
			s.unitPrefixes[uKey] = make(map[string]bool)
		}
		s.unitPrefixes[uKey][pKey] = true
	}

	return nil
}

// Clone creates a deep copy of the current System.
func (s *System) Clone() *System {
	// 1. Copy Config
	newSys := NewSystem(s.Config)

	// 2. Copy Units
	for k, u := range s.units {
		newSys.units[k] = u
	}

	// 3. Copy Prefixes
	if len(s.prefixes) > 0 {
		newSys.prefixes = make([]Prefix, len(s.prefixes))
		copy(newSys.prefixes, s.prefixes)
	}

	// 4. Copy Bindings (Deep Copy)
	for uKey, pSet := range s.unitPrefixes {
		newSet := make(map[string]bool)
		for pKey, allowed := range pSet {
			newSet[pKey] = allowed
		}
		newSys.unitPrefixes[uKey] = newSet
	}

	return newSys
}

// OverwritePrefix updates the scale of an existing prefix.
func (s *System) OverwritePrefix(symbol string, newScale float64) error {
	pKey := s.normalizeKey(symbol)

	for i, p := range s.prefixes {
		if p.Symbol == pKey {
			// Update scale directly
			s.prefixes[i].Scale = newScale
			return nil
		}
	}
	return fmt.Errorf("prefix %s not found in system, use AddPrefix instead", symbol)
}

// Resolve attempts to resolve a symbol into a Unit and a scaling factor.
func (s *System) Resolve(symbol string) (Unit, float64, bool) {
	lookupSymbol := s.normalizeKey(symbol)

	// 1. Exact Match Priority
	if u, ok := s.units[lookupSymbol]; ok {
		return u, 1.0, true
	}

	// 2. Prefix + Unit Match
	for _, p := range s.prefixes {
		pLen := len(p.Symbol)
		if len(lookupSymbol) > pLen && lookupSymbol[:pLen] == p.Symbol {
			baseSymbol := lookupSymbol[pLen:]

			// Check if the remainder is a valid unit
			if u, ok := s.units[baseSymbol]; ok {
				// Check if the prefix is allowed for this unit (Whitelist check)
				allowedPrefixes, hasList := s.unitPrefixes[baseSymbol]
				if hasList && allowedPrefixes[p.Symbol] {
					return u, p.Scale, true
				}
			}
		}
	}

	return Unit{}, 0, false
}
