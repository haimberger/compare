package basic

// ExactEqualer compares values exactly.
type ExactEqualer struct{}

// Bool compares two Boolean values exactly.
func (tc ExactEqualer) Bool(a, b bool) bool {
	return a == b
}

// Int64 compares two integer values exactly.
func (tc ExactEqualer) Int64(a, b int64) bool {
	return a == b
}

// Float64 compares two floating-point numbers exactly.
func (tc ExactEqualer) Float64(a, b float64) bool {
	return a == b
}

// String compares two string values exactly.
func (tc ExactEqualer) String(a, b string) bool {
	return a == b
}
