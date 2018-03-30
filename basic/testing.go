package basic

import "testing"

// TestExactBool tests if an equaler compares two Boolean values exactly.
func TestExactBool(t *testing.T, e Equaler) {
	type testCase struct {
		a        bool
		b        bool
		expected bool
	}
	tcs := []testCase{
		{a: false, b: false, expected: true},
		{a: false, b: true, expected: false},
		{a: true, b: false, expected: false},
		{a: true, b: true, expected: true},
	}

	for _, tc := range tcs {
		if actual := e.Bool(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}

// TestExactInt64 tests if an equaler compares two integer values exactly.
func TestExactInt64(t *testing.T, e Equaler) {
	type testCase struct {
		a        int64
		b        int64
		expected bool
	}
	tcs := []testCase{
		{a: 0, b: 0, expected: true},
		{a: 1, b: 1, expected: true},
		{a: 0, b: 1, expected: false},
		{a: 1, b: 0, expected: false},
		{a: 1, b: 2, expected: false},
		{a: 2, b: 1, expected: false},
		{a: 0, b: -1, expected: false},
		{a: -1, b: 0, expected: false},
		{a: -1, b: -1, expected: true},
		{a: -1, b: -2, expected: false},
		{a: -2, b: -1, expected: false},
		{a: -1, b: 1, expected: false},
		{a: 1, b: -1, expected: false},
	}

	for _, tc := range tcs {
		if actual := e.Int64(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}

// TestExactFloat64 tests if an equaler compares two floating-point values exactly.
func TestExactFloat64(t *testing.T, e Equaler) {
	type testCase struct {
		a        float64
		b        float64
		expected bool
	}
	tcs := []testCase{
		{a: 0, b: 0, expected: true},
		{a: 0, b: 0.1, expected: false},
		{a: 0.1, b: 0, expected: false},
		{a: 0.1, b: 0.1, expected: true},
		{a: 0.1, b: 0.2, expected: false},
		{a: 0.2, b: 0.1, expected: false},
		{a: 0, b: -0.1, expected: false},
		{a: -0.1, b: 0, expected: false},
		{a: -0.1, b: -0.1, expected: true},
		{a: -0.1, b: -0.2, expected: false},
		{a: -0.2, b: -0.1, expected: false},
		{a: -0.1, b: 0.1, expected: false},
		{a: 0.1, b: -0.1, expected: false},
	}

	for _, tc := range tcs {
		if actual := e.Float64(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}

// TestExactString tests if an equaler compares two strings exactly.
func TestExactString(t *testing.T, e Equaler) {
	type testCase struct {
		a        string
		b        string
		expected bool
	}
	tcs := []testCase{
		{a: "", b: "", expected: true},
		{a: "", b: "foo", expected: false},
		{a: "foo", b: "", expected: false},
		{a: "foo", b: "foo", expected: true},
		{a: "foo", b: "bar", expected: false},
		{a: "bar", b: "foo", expected: false},
	}

	for _, tc := range tcs {
		if actual := e.String(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}
