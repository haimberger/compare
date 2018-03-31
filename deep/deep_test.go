package deep

import (
	"math"
	"testing"
	"time"

	"github.com/haimberger/compare/basic"
)

type foo struct {
	x int
	y float32
}

type notFoo foo

func TestEqual_Exact(t *testing.T) {
	type testCase struct {
		a        interface{}
		b        interface{}
		expected bool
		err      string
	}
	// test cases copied more or less directly from https://golang.org/src/reflect/all_test.go
	tcs := []testCase{
		// Equalities
		{nil, nil, true, ""},
		{true, true, true, ""},
		{true, false, false, ""},
		{1, 1, true, ""},
		{int32(1), int32(1), true, ""},
		{0.5, 0.5, true, ""},
		{float32(0.5), float32(0.5), true, ""},
		{"hello", "hello", true, ""},
		{make([]int, 10), make([]int, 10), true, ""},
		{&[3]int{1, 2, 3}, &[3]int{1, 2, 3}, true, ""},
		{foo{1, 0.5}, foo{1, 0.5}, true, ""},
		{error(nil), error(nil), true, ""},
		{map[int]string{1: "one", 2: "two"}, map[int]string{2: "two", 1: "one"}, true, ""},
		{map[int]*bool{1: nil}, map[int]*bool{1: nil}, true, ""},

		// Inequalities
		{1, 2, false, ""},
		{int32(1), int32(2), false, ""},
		{0.5, 0.6, false, ""},
		{float32(0.5), float32(0.6), false, ""},
		{"hello", "hey", false, ""},
		{make([]int, 10), make([]int, 11), false, ""},
		{&[3]int{1, 2, 3}, &[3]int{1, 2, 4}, false, ""},
		{foo{1, 0.5}, foo{1, 0.6}, false, ""},
		{foo{1, 0}, foo{2, 0}, false, ""},
		{map[int]string{1: "one", 3: "two"}, map[int]string{2: "two", 1: "one"}, false, ""},
		{map[int]string{1: "one", 2: "txo"}, map[int]string{2: "two", 1: "one"}, false, ""},
		{map[int]string{1: "one"}, map[int]string{2: "two", 1: "one"}, false, ""},
		{map[int]string{2: "two", 1: "one"}, map[int]string{1: "one"}, false, ""},
		{nil, 1, false, ""},
		{1, nil, false, ""},
		{[][]int{{1}}, [][]int{{2}}, false, ""},
		{math.NaN(), math.NaN(), false, ""},
		{&[1]float64{math.NaN()}, &[1]float64{math.NaN()}, false, ""},
		{[]float64{math.NaN()}, []float64{math.NaN()}, false, ""},
		{map[float64]float64{math.NaN(): 1}, map[float64]float64{1: 2}, false, ""},

		// Nil vs empty: not the same.
		{[]int{}, []int(nil), false, ""},
		{[]int{}, []int{}, true, ""},
		{[]int(nil), []int(nil), true, ""},
		{map[int]int{}, map[int]int(nil), false, ""},
		{map[int]int{}, map[int]int{}, true, ""},
		{map[int]int(nil), map[int]int(nil), true, ""},

		// Mismatched types
		{1, 1.0, false, ""},
		{int32(1), int64(1), false, ""},
		{0.5, "hello", false, ""},
		{[]int{1, 2, 3}, [3]int{1, 2, 3}, false, ""},
		{&[3]interface{}{1, 2, 4}, &[3]interface{}{1, 2, "s"}, false, ""},
		{foo{1, 0.5}, notFoo{1, 0.5}, false, ""},
		{map[uint]string{1: "one", 2: "two"}, map[int]string{2: "two", 1: "one"}, false, ""},

		// Unsupported
		{uint(1), uint(1), false, "type uint not supported"},
		{2i, 2i, false, "type complex128 not supported"},
		{func() {}, func() {}, false, "type func() not supported"},

		// Empty struct (type self struct{})
		// TODO: figure out why these should be equal (they look unequal to me)
		//{&[1]float64{math.NaN()}, self{}, true, ""},
		//{[]float64{math.NaN()}, self{}, true, ""},
		//{map[float64]float64{math.NaN(): 1}, self{}, true, ""},
	}
	e := Equaler{Basic: basic.ExactEqualer{}}
	for _, tc := range tcs {
		actual, err := e.Equal(tc.a, tc.b)
		if err != nil {
			if tc.err == "" {
				t.Errorf("[%v == %v] %v", tc.a, tc.b, err)
			} else if err.Error() != tc.err {
				t.Errorf("[%v == %v] expected error %v; got %v", tc.a, tc.b, tc.err, err)
			}
		} else if tc.err != "" {
			t.Errorf("[%v == %v] expected error %v; got nothing", tc.a, tc.b, tc.err)
		} else if actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}

func TestEqual_Tolerant(t *testing.T) {
	type testCase struct {
		a        interface{}
		b        interface{}
		expected bool
	}
	tcs := []testCase{
		{a: 0.1, b: 0.151, expected: false},
		{
			a:        []float64{0.1, 0.1, 0.1, 0.1, 0.1},
			b:        []float64{0.05, 0.1, 0.10, 0.14, 0.15},
			expected: true,
		},
		// The following test fails because the tolerance is specified as a float64 value,
		// but the values are specified as float32 values, which are less precise.
		// For example, the float32 value 0.05 is actually 0.05000000074505806.
		// TODO: specify separate tolerances for float32 and float64?
		//{a: []float32{0.05}, b: []float32{0.1}, expected: true},
	}
	e := Equaler{Basic: basic.TolerantEqualer{Tolerance: 0.05}}
	for _, tc := range tcs {
		actual, err := e.Equal(tc.a, tc.b)
		if err != nil {
			t.Errorf("[%v == %v] %v", tc.a, tc.b, err)
		} else if actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}

func TestEqual_Time(t *testing.T) {
	type testCase struct {
		a        interface{}
		b        interface{}
		expected bool
	}
	tcs := []testCase{
		{
			a:        "2018-03-30T16:41:11.509Z",
			b:        "2018-03-30T16:41:12.510Z",
			expected: false,
		},
		{
			a:        []string{"2018-03-30T16:41:11.509Z", "2018-03-30T16:41:12.5Z"},
			b:        []string{"2018-03-30T16:41:11.509Z", "2018-03-30T16:41:12.509Z"},
			expected: true,
		},
	}
	tolerance, err := time.ParseDuration("1s")
	if err != nil {
		t.Fatal(err)
	}
	e := Equaler{Basic: basic.TimeEqualer{Layout: time.RFC3339Nano, Tolerance: tolerance}}
	for _, tc := range tcs {
		actual, err := e.Equal(tc.a, tc.b)
		if err != nil {
			t.Errorf("[%v == %v] %v", tc.a, tc.b, err)
		} else if actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}
