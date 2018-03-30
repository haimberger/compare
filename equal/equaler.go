package equal

// BasicEqualer specifies methods for comparing values for basic types.
type BasicEqualer interface {
	Bool(bool, bool) bool
	Int(int, int) bool
	Float64(float64, float64) bool
	String(string, string) bool
	// TODO: add functions for remaining basic types
}
