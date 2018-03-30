package basic

import (
	"math"
)

// TolerantEqualer allows some tolerance when comparing numeric values.
type TolerantEqualer struct {
	Tolerance float64
}

// Bool compares two Boolean values exactly.
func (tc TolerantEqualer) Bool(a, b bool) bool {
	return a == b
}

// Int64 compares two integer values within a tolerance. For example, if the
// tolerance is 3.7, then 8 and 11 are considered equal, but 8 and 12 are not.
func (tc TolerantEqualer) Int64(a, b int64) bool {
	return math.Abs(float64(a-b)) <= tc.Tolerance
}

// Float64 compares two floating-point numbers within a tolerance. For example,
// if the tolerance is 0.5, then 4.07 and 4.57 are considered equal, but 4.07
// and 4.571 are not.
func (tc TolerantEqualer) Float64(a, b float64) bool {
	return math.Abs(a-b) <= tc.Tolerance
}

// String compares two string values exactly.
func (tc TolerantEqualer) String(a, b string) bool {
	return a == b
}
