// Package basic provides functionality for comparing values of basic types.
package basic

// BasicEqualer implements functions for determining if values of basic types are equal.
type BasicEqualer interface {
	Bool(bool, bool) bool
	Int(int, int) bool
	Float64(float64, float64) bool
	String(string, string) bool
	// TODO: add functions for remaining basic types
}
