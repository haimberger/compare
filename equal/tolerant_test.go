package equal

import "testing"

func TestBool_Tolerant(t *testing.T) {
	// Boolean values should be compared exactly, regardless of the tolerance
	TestExactBool(t, TolerantEqualer{Tolerance: 2})
}

func TestInt_Tolerant(t *testing.T) {
	// if no tolerance is specified, values should be compared exactly
	TestExactInt(t, TolerantEqualer{})

	// if tolerance is 0, values should be compared exactly
	TestExactInt(t, TolerantEqualer{Tolerance: 0})

	// if a tolerance is set, values within the tolerance should be considered equal
	e := TolerantEqualer{Tolerance: 2}
	type testCase struct {
		a        int
		b        int
		expected bool
	}
	tcs := []testCase{
		{a: 0, b: 0, expected: true},
		{a: 0, b: 2, expected: true},
		{a: 2, b: 0, expected: true},
		{a: 0, b: 3, expected: false},
		{a: 3, b: 0, expected: false},
		{a: 3, b: 3, expected: true},
		{a: 3, b: 5, expected: true},
		{a: 5, b: 3, expected: true},
		{a: 3, b: 6, expected: false},
		{a: 6, b: 3, expected: false},
		{a: 0, b: -2, expected: true},
		{a: -2, b: 0, expected: true},
		{a: 0, b: -3, expected: false},
		{a: -3, b: 0, expected: false},
		{a: -3, b: -3, expected: true},
		{a: -3, b: -5, expected: true},
		{a: -5, b: -3, expected: true},
		{a: -3, b: -6, expected: false},
		{a: -6, b: -3, expected: false},
		{a: -3, b: 3, expected: false},
		{a: 3, b: -3, expected: false},
	}
	for _, tc := range tcs {
		if actual := e.Int(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}

func TestFloat64_Tolerant(t *testing.T) {
	// if no tolerance is specified, values should be compared exactly
	TestExactFloat64(t, TolerantEqualer{})

	// if tolerance is 0, values should be compared exactly
	TestExactFloat64(t, TolerantEqualer{Tolerance: 0})

	// if a tolerance is set, values within the tolerance should be considered equal
	e := TolerantEqualer{Tolerance: 0.05}
	type testCase struct {
		a        float64
		b        float64
		expected bool
	}
	tcs := []testCase{
		{a: 0, b: 0, expected: true},
		{a: 0, b: 0.05, expected: true},
		{a: 0.05, b: 0, expected: true},
		{a: 0, b: 0.1, expected: false},
		{a: 0.1, b: 0, expected: false},
		{a: 0.1, b: 0.1, expected: true},
		{a: 0.1, b: 0.15, expected: true},
		{a: 0.15, b: 0.1, expected: true},
		{a: 0.1, b: 0.2, expected: false},
		{a: 0.2, b: 0.1, expected: false},
		{a: 0, b: -0.05, expected: true},
		{a: -0.05, b: 0, expected: true},
		{a: 0, b: -0.1, expected: false},
		{a: -0.1, b: 0, expected: false},
		{a: -0.1, b: -0.1, expected: true},
		{a: -0.1, b: -0.15, expected: true},
		{a: -0.15, b: -0.1, expected: true},
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

func TestString_Tolerant(t *testing.T) {
	// string values should be compared exactly, regardless of the tolerance
	TestExactString(t, TolerantEqualer{Tolerance: 2})
}
