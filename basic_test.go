package compare

import (
	"regexp"
	"testing"
	"time"
)

func TestTolerantBasicEqualer_Bool(t *testing.T) {
	type testCase struct {
		a        bool
		b        bool
		expected bool
	}
	tcs := []testCase{
		{false, false, true},
		{false, true, false},
		{true, false, false},
		{true, true, true},
	}
	e := TolerantBasicEqualer{}
	for _, tc := range tcs {
		if actual := e.Bool(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}

func TestTolerantBasicEqualer_Int64(t *testing.T) {
	type testCase struct {
		a        int64
		b        int64
		expected bool
	}
	tcs := []testCase{
		{0, 0, true},
		{1, 1, true},
		{0, 1, false},
		{1, 0, false},
		{1, 2, false},
		{2, 1, false},
		{0, -1, false},
		{-1, 0, false},
		{-1, -1, true},
		{-1, -2, false},
		{-2, -1, false},
		{-1, 1, false},
		{1, -1, false},
	}
	e := TolerantBasicEqualer{}
	for _, tc := range tcs {
		if actual := e.Int64(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}

func TestTolerantBasicEqualer_Uint64(t *testing.T) {
	type testCase struct {
		a        uint64
		b        uint64
		expected bool
	}
	tcs := []testCase{
		{0, 0, true},
		{1, 1, true},
		{0, 1, false},
		{1, 0, false},
		{1, 2, false},
		{2, 1, false},
	}
	e := TolerantBasicEqualer{}
	for _, tc := range tcs {
		if actual := e.Uint64(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}

func TestTolerantBasicEqualer_Float64(t *testing.T) {
	type testCase struct {
		a        float64
		b        float64
		expected bool
	}

	exact := []testCase{
		{0, 0, true},
		{0, 0.1, false},
		{0.1, 0, false},
		{0.1, 0.1, true},
		{0.1, 0.2, false},
		{0.2, 0.1, false},
		{0, -0.1, false},
		{-0.1, 0, false},
		{-0.1, -0.1, true},
		{-0.1, -0.2, false},
		{-0.2, -0.1, false},
		{-0.1, 0.1, false},
		{0.1, -0.1, false},
	}

	approximate := []testCase{
		{0, 0, true},
		{0, 0.05, true},
		{0.05, 0, true},
		{0, 0.1, false},
		{0.1, 0, false},
		{0.1, 0.1, true},
		{0.1, 0.15, true},
		{0.15, 0.1, true},
		{0.1, 0.2, false},
		{0.2, 0.1, false},
		{0, -0.05, true},
		{-0.05, 0, true},
		{0, -0.1, false},
		{-0.1, 0, false},
		{-0.1, -0.1, true},
		{-0.1, -0.15, true},
		{-0.15, -0.1, true},
		{-0.1, -0.2, false},
		{-0.2, -0.1, false},
		{-0.1, 0.1, false},
		{0.1, -0.1, false},
	}

	// if no tolerance is specified, values should be compared exactly
	e := TolerantBasicEqualer{}
	for _, tc := range exact {
		if actual := e.Float64(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}

	// if tolerance is 0, values should be compared exactly
	e = TolerantBasicEqualer{Float64Tolerance: 0}
	for _, tc := range exact {
		if actual := e.Float64(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}

	// if a tolerance is set, values within the tolerance should be considered equal
	e = TolerantBasicEqualer{Float64Tolerance: 0.05}
	for _, tc := range approximate {
		if actual := e.Float64(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}

func TestTolerantBasicEqualer_Complex128(t *testing.T) {
	type testCase struct {
		a        complex128
		b        complex128
		expected bool
	}
	tcs := []testCase{
		{0, 0, true},
		{1i, 1i, true},
		{0, 1i, false},
		{1i, 0, false},
		{1i, 2i, false},
		{2i, 1i, false},
		{0, -1i, false},
		{-1i, 0, false},
		{-1i, -1i, true},
		{-1i, -2i, false},
		{-2i, -1i, false},
		{-1i, 1i, false},
		{1i, -1i, false},
		{1 + 2i, 2i + 1, true},
		{1 + 2i, 2 + 2i, false},
	}
	e := TolerantBasicEqualer{}
	for _, tc := range tcs {
		if actual := e.Complex128(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}

func TestTolerantBasicEqualer_String(t *testing.T) {
	// if a tolerance is specified, valid dates within the tolerance should be considered equal
	type testCase struct {
		a        string
		b        string
		expected bool
	}

	exact := []testCase{
		{"", "", true},
		{"", "foo", false},
		{"foo", "", false},
		{"foo", "foo", true},
		{"foo", "bar", false},
		{"bar", "foo", false},
	}

	approximate := []testCase{
		{"2018-03-30T16:41:11.509Z", "2018-03-30T16:41:11.509Z", true},
		{"2018-03-30T16:41:11.509Z", "2018-03-30T16:41:12.509Z", true},
		{"2018-03-30T16:41:12.509Z", "2018-03-30T16:41:11.509Z", true},
		{"2018-03-30T16:41:11.509Z", "2018-03-30T16:41:12.510Z", false},
		{"2018-03-30T16:41:12.510Z", "2018-03-30T16:41:11.509Z", false},
		{"2018-03-30T16:41:11.509Z", "foo", false},
		{"foo", "2018-03-30T16:41:11.509Z", false},
		{"2018-03-30T16:41:11.509Z", "30-Mar-18 16:41:11.509Z", false},
		{"30-Mar-18 16:41:11.509Z", "2018-03-30T16:41:11.509Z", false},
		{"foo_1_1", "foo_1_2", true},
		{"foo_1_1", "foo_2_1", false},
		{"foo", "foo", true},
		{"foo", "bar", false},
	}

	// if no special options are specified, values should be compared exactly
	e := TolerantBasicEqualer{}
	for _, tc := range exact {
		if actual := e.String(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}

	re, err := regexp.Compile("_[^_]*$") // ignore everything after last underscore
	if err != nil {
		t.Fatal(err)
	}
	tolerance, err := time.ParseDuration("1s")
	if err != nil {
		t.Fatal(err)
	}
	e = TolerantBasicEqualer{
		StringTransformer: SubstringDeleter{Regexp: re},
		TimeLayout:        time.RFC3339Nano,
		TimeTolerance:     tolerance,
	}
	for _, tc := range approximate {
		if actual := e.String(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}
