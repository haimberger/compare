package equal

import (
	"math"
	"time"
)

// TimeEqualer allows some tolerance when comparing dates.
type TimeEqualer struct {
	Layout    string
	Tolerance time.Duration
}

// Bool compares two Boolean values exactly.
func (tc TimeEqualer) Bool(a, b bool) bool {
	return a == b
}

// Int compares two integer values exactly.
func (tc TimeEqualer) Int(a, b int) bool {
	return a == b
}

// Float64 compares two floating-point numbers exactly.
func (tc TimeEqualer) Float64(a, b float64) bool {
	return a == b
}

// String compares two string values representing times within a tolerance.
// For example, if the tolerance is a second, then "2018-03-30T14:36:09.778" and
// "2018-03-30T14:36:10.778" are considered equal, but "2018-03-30T14:36:09.778"
// and "2018-03-30T14:36:10.779" are not. If at least one of the strings can't
// be parsed as a time, the strings are compared exactly.
func (tc TimeEqualer) String(a, b string) bool {
	ta, erra := time.Parse(tc.Layout, a)
	tb, errb := time.Parse(tc.Layout, b)
	if erra != nil || errb != nil {
		// TODO: log errors (at level INFO or DEBUG)?
		return a == b
	}
	diff := math.Abs(float64(ta.Sub(tb).Nanoseconds()))
	tol := float64(tc.Tolerance.Nanoseconds())
	return diff <= tol
}
