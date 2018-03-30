// Package basic provides functionality for comparing values of basic types.
package basic

// Equaler implements functions for determining if values of basic types are equal.
type Equaler interface {
	Bool(bool, bool) bool
	Int64(int64, int64) bool
	Float64(float64, float64) bool
	String(string, string) bool
	// TODO: add functions for remaining basic types
}
