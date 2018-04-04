package compare

import (
	"math"
	"testing"
	"time"
)

type foo struct {
	x int
	y float32
}

type notFoo foo

type self struct{}

func TestDeepEqualer_Equal_exact(t *testing.T) {
	type testCase struct {
		a        interface{}
		b        interface{}
		expected bool
		err      string
	}
	c := make(chan int)
	// test cases copied more or less directly from https://golang.org/src/reflect/all_test.go
	tcs := []testCase{
		// Equalities
		{nil, nil, true, ""},
		{true, true, true, ""},
		{true, false, false, ""},
		{1, 1, true, ""},
		{int32(1), int32(1), true, ""},
		{uint(1), uint(1), true, ""},
		{uintptr(1), uintptr(1), true, ""},
		{0.5, 0.5, true, ""},
		{float32(0.5), float32(0.5), true, ""},
		{2i, 2i, true, ""},
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
		{uint(1), uint(2), false, ""},
		{0.5, 0.6, false, ""},
		{float32(0.5), float32(0.6), false, ""},
		{1i, 2i, false, ""},
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
		{&[1]float64{math.NaN()}, self{}, true, ""},
		{[]float64{math.NaN()}, []float64{math.NaN()}, false, ""},
		{[]float64{math.NaN()}, self{}, true, ""},
		{map[float64]float64{math.NaN(): 1}, map[float64]float64{1: 2}, false, ""},
		{map[float64]float64{math.NaN(): 1}, self{}, true, ""},

		// Nil vs empty vs zero: not the same.
		{[]int{}, []int(nil), false, ""},
		{[]int{}, []int{}, true, ""},
		{[]int(nil), []int(nil), true, ""},
		{map[int]int{}, map[int]int(nil), false, ""},
		{map[int]int{}, map[int]int{}, true, ""},
		{map[int]int(nil), map[int]int(nil), true, ""},
		{[]interface{}{nil}, []interface{}{nil}, true, ""},
		{[]interface{}{""}, []interface{}{nil}, false, ""},
		{[]interface{}{}, []interface{}{nil}, false, ""},

		// Mismatched types
		{1, 1.0, false, ""},
		{int32(1), int64(1), false, ""},
		{0.5, "hello", false, ""},
		{[]int{1, 2, 3}, [3]int{1, 2, 3}, false, ""},
		{&[3]interface{}{1, 2, 4}, &[3]interface{}{1, 2, "s"}, false, ""},
		{foo{1, 0.5}, notFoo{1, 0.5}, false, ""},
		{map[uint]string{1: "one", 2: "two"}, map[int]string{2: "two", 1: "one"}, false, ""},

		// Unsupported
		{func() {}, func() {}, false, "type func() not supported"},
		{c, c, false, "type chan int not supported"},
	}
	// since we specify no tolerances, the equaler will compare values exactly
	e := DeepEqualer{Basic: TolerantBasicEqualer{}}
	for _, tc := range tcs {
		if tc.b == (self{}) {
			tc.b = tc.a
		}
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

func TestDeepEqualer_Equal_tolerant(t *testing.T) {
	type testCase struct {
		a        interface{}
		b        interface{}
		expected bool
	}
	tcs := []testCase{
		{0.1, 0.151, false},
		{[]float64{0.1, 0.1, 0.1, 0.1, 0.1}, []float64{0.05, 0.1, 0.10, 0.14, 0.15}, true},
		// The following test fails because the tolerance is specified as a float64 value,
		// but the values are specified as float32 values, which are less precise.
		// For example, the float32 value 0.05 is actually 0.05000000074505806.
		// TODO: specify separate tolerances for float32 and float64?
		//{[]float32{0.05}, []float32{0.1}, true},
		{"foo_1_1", "foo_2_1", false},
		{[]string{"foo_1_1", "foo_1_1"}, []string{"foo_1_2", "foo_1_abc"}, true},
		{"2018-03-30T16:41:11.509Z", "2018-03-30T16:41:12.510Z", false},
		{
			[]string{"2018-03-30T16:41:11.509Z", "2018-03-30T16:41:12.5Z"},
			[]string{"2018-03-30T16:41:11.509Z", "2018-03-30T16:41:12.509Z"},
			true,
		},
	}
	sd, err := MkSubstringDeleter("_[^_]*$") // ignore everything after last underscore
	if err != nil {
		t.Fatal(err)
	}
	tolerance, err := time.ParseDuration("1s")
	if err != nil {
		t.Fatal(err)
	}
	e := DeepEqualer{Basic: TolerantBasicEqualer{
		Float64Tolerance:  0.05,
		StringTransformer: sd,
		TimeLayout:        time.RFC3339Nano,
		TimeTolerance:     tolerance,
	}}
	for _, tc := range tcs {
		actual, err := e.Equal(tc.a, tc.b)
		if err != nil {
			t.Errorf("[%v == %v] %v", tc.a, tc.b, err)
		} else if actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}
