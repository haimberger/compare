package basic

import (
	"testing"
	"time"
)

func TestBool_Time(t *testing.T) {
	// Boolean values should be compared exactly, regardless of the tolerance
	TestExactBool(t, TimeEqualer{})
}

func TestInt_Time(t *testing.T) {
	// integer values should be compared exactly, regardless of the tolerance
	TestExactInt64(t, TimeEqualer{})
}

func TestFloat64_Time(t *testing.T) {
	// floating-point values should be compared exactly, regardless of the tolerance
	TestExactFloat64(t, TimeEqualer{})
}

func TestString_Time(t *testing.T) {
	// if a tolerance is specified, valid dates within the tolerance should be considered equal
	type testCase struct {
		a        string
		b        string
		expected bool
	}
	tcs := []testCase{
		{a: "2018-03-30T16:41:11.509Z", b: "2018-03-30T16:41:11.509Z", expected: true},
		{a: "2018-03-30T16:41:11.509Z", b: "2018-03-30T16:41:12.509Z", expected: true},
		{a: "2018-03-30T16:41:12.509Z", b: "2018-03-30T16:41:11.509Z", expected: true},
		{a: "2018-03-30T16:41:11.509Z", b: "2018-03-30T16:41:12.510Z", expected: false},
		{a: "2018-03-30T16:41:12.510Z", b: "2018-03-30T16:41:11.509Z", expected: false},
		{a: "2018-03-30T16:41:11.509Z", b: "foo", expected: false},
		{a: "foo", b: "2018-03-30T16:41:11.509Z", expected: false},
		{a: "2018-03-30T16:41:11.509Z", b: "30-Mar-18 16:41:11.509Z", expected: false},
		{a: "30-Mar-18 16:41:11.509Z", b: "2018-03-30T16:41:11.509Z", expected: false},
		{a: "foo", b: "foo", expected: true},
		{a: "foo", b: "bar", expected: false},
	}
	tolerance, err := time.ParseDuration("1s")
	if err != nil {
		t.Fatal(err)
	}
	e := TimeEqualer{Layout: time.RFC3339Nano, Tolerance: tolerance}
	for _, tc := range tcs {
		if actual := e.String(tc.a, tc.b); actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}
