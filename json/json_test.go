package json

import (
	"testing"
	"time"

	"github.com/haimberger/compare/basic"
)

func TestEqual_Exact(t *testing.T) {
	type testCase struct {
		a        string
		b        string
		expected bool
		err      string
	}
	tcs := []testCase{
		{a: ``, b: ``, err: "unexpected end of JSON input"},
		{a: `""`, b: ``, err: "unexpected end of JSON input"},
		{a: `undefined`, b: `undefined`, err: "invalid character 'u' looking for beginning of value"},
		{a: `null`, b: `null`, expected: true},
		{a: `null`, b: `false`, expected: false},
		{a: `null`, b: `""`, expected: false},
		{a: `false`, b: `false`, expected: true},
		{a: `false`, b: `true`, expected: false},
		{a: `1`, b: `1`, expected: true},
		{a: `1`, b: `2`, expected: false},
		{a: `0.1`, b: `0.1`, expected: true},
		{a: `0.1`, b: `0.2`, expected: false},
		{a: `"foo"`, b: `"foo"`, expected: true},
		{a: `"foo"`, b: `"bar"`, expected: false},
		{a: `[false, 1, 0.1, "foo"]`, b: `[false,1,0.100,"foo"]`, expected: true},
		{a: `[false, 1, 0.1, "foo"]`, b: `[false, 1, 0.2, "foo"]`, expected: false},
		{a: `[false, 1, 0.1, "foo"]`, b: `[false, 0.1, 1, "foo"]`, expected: false},
		{a: `[false, 1, 0.1]`, b: `[false, 1, 0.1, "foo"]`, expected: false},
		{
			a:        `{"a": false, "b": [1, {"c": 0.1, "d": "foo"}]}`,
			b:        `{"b": [1, {"d": "foo", "c": 0.1}],"a":false}`,
			expected: true,
		},
		{
			a:        `{"b": [1, {"c": 0.1, "d": "foo"}]}`,
			b:        `{"a": false, "b": [1, {"c": 0.1, "d": "foo"}]}`,
			expected: false,
		},
		{
			a:        `{"a": false, "b": [1, {"c": 0.1, "d": "foo"}]}`,
			b:        `{"a": false, "b": [1, {"d": 0.1, "c": "foo"}]}`,
			expected: false,
		},
	}
	be := basic.ExactEqualer{}
	for _, tc := range tcs {
		actual, err := Equal([]byte(tc.a), []byte(tc.b), be)
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
		a        string
		b        string
		expected bool
		err      string
	}
	tcs := []testCase{
		{a: `0.1`, b: `0.151`, expected: false},
		{a: `[0.1,0.1,0.1,0.1,0.1]`, b: `[0.05,0.1,0.10,0.14,0.15]`, expected: true},
	}
	be := basic.TolerantEqualer{Tolerance: 0.05}
	for _, tc := range tcs {
		actual, err := Equal([]byte(tc.a), []byte(tc.b), be)
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

func TestEqual_Time(t *testing.T) {
	type testCase struct {
		a        string
		b        string
		expected bool
		err      string
	}
	tcs := []testCase{
		{
			a:        `"2018-03-30T16:41:11.509Z"`,
			b:        `"2018-03-30T16:41:12.510Z"`,
			expected: false,
		},
		{
			a:        `["2018-03-30T16:41:11.509Z","2018-03-30T16:41:12.5Z"]`,
			b:        `["2018-03-30T16:41:11.509Z","2018-03-30T16:41:12.509Z"]`,
			expected: true,
		},
	}
	tolerance, err := time.ParseDuration("1s")
	if err != nil {
		t.Fatal(err)
	}
	be := basic.TimeEqualer{Layout: time.RFC3339Nano, Tolerance: tolerance}
	for _, tc := range tcs {
		actual, err := Equal([]byte(tc.a), []byte(tc.b), be)
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
