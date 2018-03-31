package compare

import (
	"math"
	"regexp"
	"time"
)

// BasicEqualer provides functions for determining if values of basic types are equal.
type BasicEqualer interface {
	// Bool determines if two Boolean values are equal.
	Bool(bool, bool) bool
	// Int64 determines if two integer values are equal.
	Int64(int64, int64) bool
	// Float64 determines if two floating-point values are equal.
	Float64(float64, float64) bool
	// String determines if two strings are equal.
	String(string, string) bool
	// TODO: add functions for remaining basic types
}

// StringTransformer provides a function for transforming strings.
type StringTransformer interface {
	// Transform transforms a string into another string.
	Transform(string) string
}

// SubstringDeleter provides functionality for deleting substrings matching a
// regular expression.
type SubstringDeleter struct {
	// Regexp specifies which substrings should be deleted.
	Regexp *regexp.Regexp
}

// MkSubstringDeleter creates a new SubstringDeleter based on the specified
// regular expression.
func MkSubstringDeleter(expr string) (SubstringDeleter, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return SubstringDeleter{}, err
	}
	return SubstringDeleter{Regexp: re}, nil
}

// Transform removes all substrings matching a regular expression.
func (sd SubstringDeleter) Transform(s string) string {
	return sd.Regexp.ReplaceAllString(s, "")
}

// TolerantBasicEqualer is an example implementation of the Equaler interface.
// Rather than comparing values exactly, it allows some leeway.
type TolerantBasicEqualer struct {
	// Float64Tolerance specifies how much two floating-point values may differ
	// while still being considered equal.
	Float64Tolerance float64
	// StringTransformer specifies how string values should be transformed
	// before comparing them.
	StringTransformer StringTransformer
	// TimeLayout specifies the layout of time values.
	TimeLayout string
	// TimeTolerance specifies how much two times may differ while still being
	// considered equal.
	TimeTolerance time.Duration
}

// Bool compares two Boolean values exactly.
func (e TolerantBasicEqualer) Bool(a, b bool) bool {
	return a == b
}

// Int64 compares two integer values exactly.
func (e TolerantBasicEqualer) Int64(a, b int64) bool {
	return a == b
}

// Float64 compares two floating-point numbers within a tolerance. For example,
// if the tolerance is 0.5, then 4.07 and 4.57 are considered equal, but 4.07
// and 4.571 are not.
func (e TolerantBasicEqualer) Float64(a, b float64) bool {
	return math.Abs(a-b) <= e.Float64Tolerance
}

// String compares two string values.
//
// If a tolerance for time values is specified, and if both values represent
// valid times according to the specified layout (time.RFC3339 by default), then
// this function checks if the times are within the specified tolerance.
// For example, if the tolerance is a second, then "2018-03-30T14:36:09.778" and
// "2018-03-30T14:36:10.778" are considered equal, but "2018-03-30T14:36:09.778"
// and "2018-03-30T14:36:10.779" are not.
//
// Otherwise, if a StringTransformer is specified, it will be used to transform
// both strings before comparing them.
//
// If neither of the above applies, the strings are compared exactly.
func (e TolerantBasicEqualer) String(a, b string) bool {
	// if a tolerance for time values is specified, try comparing the strings as times
	if e.TimeTolerance.Nanoseconds() > 0 {
		ta, erra := time.Parse(e.TimeLayout, a)
		tb, errb := time.Parse(e.TimeLayout, b)
		if erra == nil && errb == nil {
			diff := math.Abs(float64(ta.Sub(tb).Nanoseconds()))
			tol := float64(e.TimeTolerance.Nanoseconds())
			return diff <= tol
		}
	}

	// try transforming strings before comparing them
	if e.StringTransformer != nil {
		ta := e.StringTransformer.Transform(a)
		tb := e.StringTransformer.Transform(b)
		return ta == tb
	}

	// if all else fails, compare the strings exactly
	return a == b
}
